package main

import (
	. "github.com/franela/goblin"
	"github.com/garyburd/redigo/redis"
	"github.com/rafaeljusto/redigomock"
	"testing"
)

func TestPeriodManager(t *testing.T) {

	g := Goblin(t)

	g.Describe("PeriodManager", func() {
		var redisPool *redis.Pool

		g.Before(func() {
			redisPool = redis.NewPool(func() (redis.Conn, error) {
				conn := redigomock.NewConn()
				var err error
				return conn, err
			}, 10)
		})

		g.It("should add each trade to the period manager", func() {

			periodHash := NewPeriodHash(redisPool, "")
			periodHash.setEnv("TESTING")

			periodManager := NewPeriodManager(redisPool, periodHash, "")
			periodManager.setEnv("TESTING")
			periodManager.setTickers([]string{"STOCK"})
			periodManager.initPeriods()

			periodManager.add(AnonymizedTrade{
				Shares: 150, Ticker: "STOCK", Price: 10.50,
			})
			periodManager.add(AnonymizedTrade{
				Shares: 10, Ticker: "STOCK", Price: 11.50,
			})
			periodManager.add(AnonymizedTrade{
				Shares: 30, Ticker: "STOCK", Price: 11.40,
			})
			periodManager.add(AnonymizedTrade{
				Shares: 30, Ticker: "STOCK", Price: 11.10,
			})
			periodManager.add(AnonymizedTrade{
				Shares: 20, Ticker: "STOCK", Price: 9.50,
			})
			periodManager.add(AnonymizedTrade{
				Shares: 30, Ticker: "STOCK", Price: 10.55,
			})

			thing := periodManager.hash.get("STOCK")

			// cant mock the time, so we'll test each property that isnt time
			g.Assert(thing.High).Equal(11.5)
			g.Assert(thing.Low).Equal(9.5)
			g.Assert(thing.Open).Equal(10.5)
			g.Assert(thing.Close).Equal(10.55)
			g.Assert(thing.Volume).Equal(270)
			g.Assert(thing.Ticker).Equal("STOCK")
		})
	})
}
