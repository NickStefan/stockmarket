package main

import (
	. "github.com/franela/goblin"
	"testing"
)

func TestPeriodManager(t *testing.T) {

	g := Goblin(t)

	g.Describe("PeriodManager", func() {

		g.It("should add each trade to the period manager", func() {

			periodManager := NewPeriodManager([]string{"STOCK"})

			periodManager.add(Trade{
				Shares: 150, Ticker: "STOCK", Price: 10.50,
			})
			periodManager.add(Trade{
				Shares: 10, Ticker: "STOCK", Price: 11.50,
			})
			periodManager.add(Trade{
				Shares: 30, Ticker: "STOCK", Price: 11.40,
			})
			periodManager.add(Trade{
				Shares: 30, Ticker: "STOCK", Price: 11.10,
			})
			periodManager.add(Trade{
				Shares: 20, Ticker: "STOCK", Price: 9.50,
			})
			periodManager.add(Trade{
				Shares: 30, Ticker: "STOCK", Price: 10.55,
			})

			// cant mock the time, so we'll test each property that isnt time
			g.Assert(periodManager.hash.get("STOCK").High).Equal(11.5)
			g.Assert(periodManager.hash.get("STOCK").Low).Equal(9.5)
			g.Assert(periodManager.hash.get("STOCK").Open).Equal(10.5)
			g.Assert(periodManager.hash.get("STOCK").Close).Equal(10.55)
			g.Assert(periodManager.hash.get("STOCK").Volume).Equal(270)
			g.Assert(periodManager.hash.get("STOCK").Ticker).Equal("STOCK")
		})
	})
}
