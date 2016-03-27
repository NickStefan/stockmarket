package main

import (
	"testing"
	. "github.com/franela/goblin"
)

func Test(t *testing.T){
	g := Goblin(t)

	g.Describe("Account", func(){

		var dataStore map[string]*Account
		
		g.BeforeEach(func(){
			dataStore = map[string]*Account{
				"Bob": &Account{
					name: "Bob",
					cash: 0,
					assets: map[string]*Asset{
						"STOCK": &Asset{ ticker: "STOCK", shares: 100, },
					},
				},
				"Tim": &Account{
					name: "Tim",
					cash: 10000,
					assets: make(map[string]*Asset),
				},
			}
		})

		g.It("should increment and decrement shares in Assets hash", func(){
			tradeBob := Trade{
				Actor: "Bob",
				Shares: 100,
				Ticker: "STOCK",
				Price: 7,
				Intent: "SELL",
				Kind: "LIMIT",
			}
			tradeTim := Trade{
				Actor: "Tim",
				Shares: 100,
				Ticker: "STOCK",
				Price: 100,
				Intent: "BUY",
				Kind: "MARKET",
			}

			processTrade(dataStore, tradeBob)
			processTrade(dataStore, tradeTim)

			g.Assert(dataStore["Tim"].assets["STOCK"].shares).Equal(0)
			g.Assert(dataStore["Bob"].assets["STOCK"].shares).Equal(100)
		})

		g.It("should update the cash balance", func(){

			tradeBob := Trade{
				Actor: "Bob",
				Shares: 100,
				Ticker: "STOCK",
				Price: 7,
				Intent: "SELL",
				Kind: "LIMIT",
			}
			tradeTim := Trade{
				Actor: "Tim",
				Shares: 100,
				Ticker: "STOCK",
				Price: 100,
				Intent: "BUY",
				Kind: "MARKET",
			}

			processTrade(dataStore, tradeBob)
			processTrade(dataStore, tradeTim)

			g.Assert(dataStore["Tim"].cash).Equal(0.00)
			g.Assert(dataStore["Bob"].cash).Equal(10000.00)
		})
		
	})
}