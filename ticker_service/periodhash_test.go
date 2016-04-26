package main

import (
	. "github.com/franela/goblin"
	"testing"
)

func TestPeriodHash(t *testing.T) {

	g := Goblin(t)

	g.Describe("PeriodHash", func() {

		g.It("should add each trade to the period cache", func() {

			periodHash := NewPeriodHash([]string{"STOCK"})

			periodHash.add(Trade{
				Shares: 150, Ticker: "STOCK", Price: 10.50,
			})
			periodHash.add(Trade{
				Shares: 10, Ticker: "STOCK", Price: 11.50,
			})
			periodHash.add(Trade{
				Shares: 30, Ticker: "STOCK", Price: 11.40,
			})
			periodHash.add(Trade{
				Shares: 30, Ticker: "STOCK", Price: 11.10,
			})
			periodHash.add(Trade{
				Shares: 20, Ticker: "STOCK", Price: 9.50,
			})
			periodHash.add(Trade{
				Shares: 30, Ticker: "STOCK", Price: 10.55,
			})

			// cant mock the time, so we'll test each property that isnt time
			g.Assert(periodHash.hash["STOCK"].High).Equal(11.5)
			g.Assert(periodHash.hash["STOCK"].Low).Equal(9.5)
			g.Assert(periodHash.hash["STOCK"].Open).Equal(10.5)
			g.Assert(periodHash.hash["STOCK"].Close).Equal(10.55)
			g.Assert(periodHash.hash["STOCK"].Volume).Equal(270)
			g.Assert(periodHash.hash["STOCK"].Ticker).Equal("STOCK")
		})
	})
}
