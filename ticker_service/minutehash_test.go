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

            g.Assert(minuteHash.hash["STOCK"]).Equal(&Minute{
                High: 11.5, Low: 9.5, 
                Open: 10.5, Close: 10.55, 
                Volume: 270, Ticker: "STOCK",
            })
		})
	})
}