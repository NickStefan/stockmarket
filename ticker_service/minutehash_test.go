package main

import (
	"testing"
	. "github.com/franela/goblin"
)

func TestMinuteHash(t *testing.T) {

	g := Goblin(t)

	g.Describe("MinuteHash", func(){

		g.It("should add each trade to the minute cache", func(){

			minuteHash := NewMinuteHash([]string{"STOCK"})

            minuteHash.add(Trade{
                Shares: 150, Ticker: "STOCK", Price: 10.50,
            })
            minuteHash.add(Trade{
                Shares: 10, Ticker: "STOCK", Price: 11.50,
            })
            minuteHash.add(Trade{
                Shares: 30, Ticker: "STOCK", Price: 11.40,
            })
            minuteHash.add(Trade{
                Shares: 30, Ticker: "STOCK", Price: 11.10,
            })
            minuteHash.add(Trade{
                Shares: 20, Ticker: "STOCK", Price: 9.50,
            })
            minuteHash.add(Trade{
                Shares: 30, Ticker: "STOCK", Price: 10.55,
            })

            // cant mock the time, so we'll test each property that isnt time
            g.Assert(minuteHash.hash["STOCK"].High).Equal(11.5)
            g.Assert(minuteHash.hash["STOCK"].Low).Equal(9.5)
            g.Assert(minuteHash.hash["STOCK"].Open).Equal(10.5)
            g.Assert(minuteHash.hash["STOCK"].Close).Equal(10.55)
            g.Assert(minuteHash.hash["STOCK"].Volume).Equal(270)
            g.Assert(minuteHash.hash["STOCK"].Ticker).Equal("STOCK")   
		})
	})
}