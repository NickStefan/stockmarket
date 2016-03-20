package main

import (
	"testing"
	. "github.com/franela/goblin"
)

func Test(t *testing.T){
	g := Goblin(t)

	g.Describe("Orders", func(){

		var orders [6]Order

		g.BeforeEach(func(){
			
			orders = [6]Order{
				BuyLimit{
					bid: 10.05, 
					BaseOrder:BaseOrder{
						actor: "Bob", timecreated: "10:00am", intent: "BUY",
						shares: 100, state: "OPEN",
					},
				},
				BuyMarket{
					BaseOrder:BaseOrder{
						actor: "Tim", timecreated: "9:58am", intent: "BUY",
						shares: 100, state: "OPEN",
					},
				},
				BuyLimit{
					bid: 10.00, 
					BaseOrder:BaseOrder{
						actor: "Bob", timecreated: "10:00am", intent: "BUY",
						shares: 100, state: "OPEN",
					},
				},
				SellMarket{
					BaseOrder:BaseOrder{
						actor: "Tim", timecreated: "9:58am", intent: "SELL",
						shares: 100, state: "OPEN",
					},
				},
				SellLimit{
					ask: 10.10, 
					BaseOrder:BaseOrder{
						actor: "Tim", timecreated: "9:58am", intent: "SELL",
						shares: 100, state: "OPEN",
					},
				},
				SellMarket{
					BaseOrder:BaseOrder{
						actor: "Tim", timecreated: "9:58am", intent: "SELL",
						shares: 100, state: "OPEN",
					},
				},
			}
		})

		g.Describe("Order Interface", func(){

			g.Describe("price method", func(){
				g.It("should equal the bid on a BuyLimit", func(){
					g.Assert(orders[0].price()).Equal(10.05)
				})

				g.It("should equal the ask on a SellLimit", func(){
					g.Assert(orders[4].price()).Equal(10.10)
				})

				g.It("should equal nearly infinity on BuyMarket", func(){
					g.Assert(orders[1].price()).Equal(1000000.00)
				})

				g.It("should equal 0 on SellMarket", func(){
					g.Assert(orders[5].price()).Equal(0.00)
				})
			})

			g.Describe("lookup method", func(){
				g.It("should equal actor + createtime", func(){
					g.Assert(orders[0].lookup()).Equal("Bob10:00am")
				})
			})

			g.Describe("getOrder method", func(){
				g.It("should provide access to the embeded order struct", func(){
					g.Assert(orders[0].getOrder().shares).Equal(100)
				})
			})
		})

		g.Describe("OrderBook", func(){

			g.It("should add orders to the correct queues and hashes", func(){
				orderBook := NewOrderBook()

				for i := 0; i < len(orders); i++ {
					orderBook.add(orders[i])
				}

				g.Assert(orderBook.buyQueue.Dequeue().Value).Equal(1000000.00)
				g.Assert(orderBook.buyQueue.Dequeue().Value).Equal(10.05)
				g.Assert(orderBook.buyQueue.Dequeue().Value).Equal(10.00)
				
				g.Assert(orderBook.sellQueue.Dequeue().Value).Equal(0.00)
				g.Assert(orderBook.sellQueue.Dequeue().Value).Equal(0.00)
				g.Assert(orderBook.sellQueue.Dequeue().Value).Equal(10.10)
			})

			g.It("should fill the highest priority orders until no more can be filled", func(){

				orderBook := NewOrderBook()

				for i :=0; i < len(orders); i++ {
					orderBook.add(orders[i])
				}

				// filling orders will dequeue filled orders,
				// so expect further down the line orders when dequeueing
				orderBook.run()
				g.Assert(orderBook.buyQueue.Dequeue().Value).Equal(10.00)
				g.Assert(orderBook.sellQueue.Dequeue().Value).Equal(10.10)
			})

			g.It("should work with repeated calls to add and run", func(){
				orderBook := NewOrderBook()

				for i :=0; i < len(orders); i++ {
					orderBook.add(orders[i])
				}

				orderBook.run()

				orderBook.add(orders[1])

				orderBook.run()

				g.Assert(orderBook.buyQueue.Dequeue().Value).Equal(10.00)
				g.Assert(orderBook.sellQueue.Dequeue() == nil).Equal(true)
			})
		})

	})
}

// BENCHMARKS
var benchOrders = [6]Order{
	BuyLimit{
		bid: 10.05, 
		BaseOrder:BaseOrder{
			actor: "Bob", timecreated: "10:00am", intent: "BUY",
			shares: 100, state: "OPEN",
		},
	},
	BuyMarket{
		BaseOrder:BaseOrder{
			actor: "Tim", timecreated: "9:58am", intent: "BUY",
			shares: 100, state: "OPEN",
		},
	},
	BuyLimit{
		bid: 10.00, 
		BaseOrder:BaseOrder{
			actor: "Bob", timecreated: "10:00am", intent: "BUY",
			shares: 100, state: "OPEN",
		},
	},
	SellMarket{
		BaseOrder:BaseOrder{
			actor: "Tim", timecreated: "9:58am", intent: "SELL",
			shares: 100, state: "OPEN",
		},
	},
	SellLimit{
		ask: 10.10, 
		BaseOrder:BaseOrder{
			actor: "Tim", timecreated: "9:58am", intent: "SELL",
			shares: 100, state: "OPEN",
		},
	},
	SellMarket{
		BaseOrder:BaseOrder{
			actor: "Tim", timecreated: "9:58am", intent: "SELL",
			shares: 100, state: "OPEN",
		},
	},
}

var result float64

func BenchmarkOrderBookRun(b *testing.B){

	orderBook := NewOrderBook()

	for i :=0; i < len(benchOrders); i++ {
		orderBook.add(benchOrders[i])
	}

	// filling orders will dequeue filled orders,
	// so expect further down the line orders when dequeueing
	orderBook.run()
	result = orderBook.buyQueue.Dequeue().Value
}