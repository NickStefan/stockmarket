package main

import (
	"testing"
	. "github.com/franela/goblin"
)

func Test(t *testing.T){
	g := Goblin(t)

	g.Describe("Order Interface", func(){

		g.Describe("price method", func(){
			g.It("should equal the bid on a BuyLimit", func(){
				buyLimit := BuyLimit{
					bid: 10.00, 
					BaseOrder:BaseOrder{
						actor: "Bob", timecreated: "10:00am", intent: "BUY",
						shares: 100, state: "OPEN",
					},
				}
				g.Assert(buyLimit.price()).Equal(10.00)
			})

			g.It("should equal the ask on a SellLimit", func(){
				sellLimit := SellLimit{
					ask: 10.05, 
					BaseOrder:BaseOrder{
						actor: "Tim", timecreated: "9:58am", intent: "SELL",
						shares: 100, state: "OPEN",
					},
				}
				g.Assert(sellLimit.price()).Equal(10.05)
			})

			g.It("should equal nearly infinity on BuyMarket", func(){
				buyMarket := BuyMarket{
					BaseOrder:BaseOrder{
						actor: "Tim", timecreated: "9:58am", intent: "SELL",
						shares: 100, state: "OPEN",
					},
				}
				g.Assert(buyMarket.price()).Equal(1000000.00)
			})

			g.It("should equal 0 on SellMarket", func(){
				sellMarket := SellMarket{
					BaseOrder:BaseOrder{
						actor: "Tim", timecreated: "9:58am", intent: "SELL",
						shares: 100, state: "OPEN",
					},
				}
				g.Assert(sellMarket.price()).Equal(0.00)
			})
		})

		g.Describe("lookup method", func(){
			g.It("should equal actor + createtime", func(){
				buyMarket := SellMarket{
					BaseOrder:BaseOrder{
						actor: "Tim", timecreated: "9:58am", intent: "SELL",
						shares: 100, state: "OPEN",
					},
				}
				g.Assert(buyMarket.lookup()).Equal("Tim9:58am")
			})
		})

		g.Describe("getOrder method", func(){
			g.It("should provide access to the embeded order struct", func(){
				buyMarket := SellMarket{
					BaseOrder:BaseOrder{
						actor: "Tim", timecreated: "9:58am", intent: "SELL",
						shares: 100, state: "OPEN",
					},
				}
				g.Assert(buyMarket.getOrder().shares).Equal(100)
			})
		})
	})

	g.Describe("OrderBook", func(){
		g.It("should add orders to the correct queues and hashes", func(){

		})

		g.It("should fill the highest priority orders until no more can be filled", func(){

		})

		g.It("should be resilient after repeated adds and fills", func(){

		})
	})
}