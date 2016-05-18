package main

import (
	"github.com/garyburd/redigo/redis"
	"github.com/rafaeljusto/redigomock"
	"testing"
	"time"
	// "fmt"
	. "github.com/franela/goblin"
)

func createDummyOrders(n int64) [7]Order {
	return [7]Order{
		Order{
			Bid:   10.05,
			Actor: "Bob", Timecreated: time.Now().Unix() + n,
			Intent: "BUY", Kind: "LIMIT", Shares: 100, State: "OPEN",
		},
		Order{
			Actor: "Tim", Timecreated: time.Now().Unix() + n,
			Intent: "BUY", Kind: "MARKET", Shares: 100, State: "OPEN",
		},
		Order{
			Bid:   10.00,
			Actor: "Gary", Timecreated: time.Now().Unix() + n,
			Intent: "BUY", Kind: "LIMIT", Shares: 100, State: "OPEN",
		},
		Order{
			Actor: "Terry", Timecreated: time.Now().Unix() + n,
			Intent: "SELL", Kind: "MARKET", Shares: 100, State: "OPEN",
		},
		Order{
			Ask:   10.10,
			Actor: "Larry", Timecreated: time.Now().Unix() + n,
			Intent: "SELL", Kind: "LIMIT", Shares: 100, State: "OPEN",
		},
		Order{
			Actor: "Sam", Timecreated: time.Now().Unix() + n,
			Intent: "SELL", Kind: "MARKET", Shares: 100, State: "OPEN",
		},
		Order{
			Bid:   10.05,
			Actor: "Sally", Timecreated: time.Now().Unix() + n,
			Intent: "BUY", Kind: "LIMIT", Shares: 80, State: "OPEN",
		},
	}
}

func Test(t *testing.T) {
	g := Goblin(t)

	g.Describe("Orders", func() {

		var orders [7]Order
		var moreOrders [7]Order

		g.BeforeEach(func() {
			orders = createDummyOrders(0)
			moreOrders = createDummyOrders(1)
		})

		g.Describe("Order", func() {

			g.Describe("price method", func() {
				g.It("should equal the bid on a BuyLimit", func() {
					g.Assert(orders[0].price()).Equal(10.05)
				})

				g.It("should equal the ask on a SellLimit", func() {
					g.Assert(orders[4].price()).Equal(10.10)
				})

				g.It("should equal nearly infinity on BuyMarket", func() {
					g.Assert(orders[1].price()).Equal(1000000.00)
				})

				g.It("should equal 0 on SellMarket", func() {
					g.Assert(orders[5].price()).Equal(0.00)
				})
			})

			g.Describe("lookup method", func() {
				g.It("should equal actor + createtime", func() {
					g.Assert(orders[0].lookup()[:3]).Equal("Bob")
				})
			})

			g.Describe("properties", func() {
				g.It("should provide access to properties", func() {
					g.Assert(orders[0].Shares).Equal(100)
				})
			})
		})

		g.Describe("OrderBook", func() {

			var redisPool *redis.Pool

			g.Before(func() {
				redisPool = redis.NewPool(func() (redis.Conn, error) {
					conn := redigomock.NewConn()
					var err error
					return conn, err
				}, 10)
			})

			g.It("should add orders to the correct queues and hashes", func() {
				orderBook := NewOrderBook(redisPool)
				orderBook.setEnv("TESTING")

				for i := 0; i < 6; i++ {
					orderBook.add(&orders[i])
				}

				g.Assert(orderBook.buyQueue.Dequeue().Value).Equal(1000000.00)
				g.Assert(orderBook.buyQueue.Dequeue().Value).Equal(10.05)
				g.Assert(orderBook.buyQueue.Dequeue().Value).Equal(10.00)

				g.Assert(orderBook.sellQueue.Dequeue().Value).Equal(0.00)
				g.Assert(orderBook.sellQueue.Dequeue().Value).Equal(0.00)
				g.Assert(orderBook.sellQueue.Dequeue().Value).Equal(10.10)
			})

			g.It("should fill the highest priority orders until no more can be filled", func() {

				orderBook := NewOrderBook(redisPool)
				orderBook.setEnv("TESTING")

				for i := 0; i < 6; i++ {
					orderBook.add(&orders[i])
				}

				// filling orders will dequeue filled orders,
				// so expect further down the line orders when dequeueing
				orderBook.run()
				g.Assert(orderBook.buyQueue.Dequeue().Value).Equal(10.00)
				g.Assert(orderBook.sellQueue.Dequeue().Value).Equal(10.10)
			})

			g.It("should work with repeated calls to add and run", func() {
				orderBook := NewOrderBook(redisPool)
				orderBook.setEnv("TESTING")

				for i := 0; i < 6; i++ {
					orderBook.add(&orders[i])
				}

				orderBook.run()

				orderBook.add(&orders[1])

				orderBook.run()

				g.Assert(orderBook.buyQueue.Dequeue().Value).Equal(10.00)
				g.Assert(orderBook.sellQueue.Dequeue() == nil).Equal(true)
			})

			g.It("should call tradeHandler with matched orders", func() {

				orderBook := NewOrderBook(redisPool)
				orderBook.setEnv("TESTING")

				orderBook.setTradeHandler(func(t Trade, o Trade) {
					g.Assert(t.Price).Equal(10.05)
					g.Assert(o.Price).Equal(10.05)
				})

				orderBook.add(&orders[0])
				orderBook.add(&orders[3])
				orderBook.run()
			})

			g.It("should partially fill orders when the share numbers dont match", func() {
				orderBook := NewOrderBook(redisPool)
				orderBook.setEnv("TESTING")

				orderBook.setTradeHandler(func(t Trade, o Trade) {
					g.Assert(t.Price).Equal(10.05)
					g.Assert(o.Price).Equal(10.05)
				})

				// set one order to be smaller than the other
				orderBook.add(&orders[6])
				orderBook.add(&orders[3])
				orderBook.run()

				var lookup = orderBook.sellQueue.Peek().Lookup
				var thing = orderBook.sellHash.get(lookup)
				g.Assert(thing.Shares).Equal(20)
			})

		})

	})
}
