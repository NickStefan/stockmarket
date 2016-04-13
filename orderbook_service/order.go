package main

import (
	"time"
 	"strconv"
)

// what makes up a stock market? competing buy and sell orders

// market buy order "i will buy 100 shares at the lowest available sell price"
// limit buy order "i will buy 100 shares at the lowest available price below ____"

// market sell order "i will sell 100 shares at the highest available price"
// limit sell order "i will sell 100 shares highest available price above _____"

type BaseOrder struct {
	Shares int `json:"shares"`
	Ticker string `json:"ticker"` // STOCK
	Actor string `json:"actor"`   // BOB
	Intent string `json:"intent"` // BUY || SELL
	Kind string `json:"kind"`			// MARKET || LIMIT
	State string `json:"state"`   // OPEN || FILLED || CANCELED
	Timecreated int64 `json:"timecreated"` // unix time
}

func (b *BaseOrder) lookup() string {
	return b.Actor + strconv.FormatInt(b.Timecreated, 10)
}

func (b *BaseOrder) getOrder() *BaseOrder {
	return b
}

func (b *BaseOrder) partialFill(price float64, newShares int) Trade {
	b.Shares = newShares
	return Trade{
		Actor: b.Actor, Shares: b.Shares - newShares,
		Price: price, Intent: b.Intent, Kind: b.Kind, Ticker: b.Ticker,
		Time: time.Now().Unix(),
	}
}

func (b *BaseOrder) fill(price float64) Trade {
	return Trade{
		Actor: b.Actor, Shares: b.Shares, Price: price,
		Intent: b.Intent, Ticker: b.Ticker, Kind: b.Kind,
		Time: time.Now().Unix(),
	}
}

type BuyLimit struct {
	Bid float64 `json:"bid"`
	*BaseOrder `json:"baseorder"`
}

type SellLimit struct {
	Ask float64 `json:"ask"`
	*BaseOrder `json:"baseorder"`
}

type BuyMarket struct {
	*BaseOrder `json:"baseorder"`
}

type SellMarket struct {
	*BaseOrder `json:"baseorder"`
}

// create a consistent interface for the different types of orders

type Order interface {
	price() float64
	lookup() string
	getOrder() *BaseOrder
	partialFill(price float64, newShares int) Trade
	fill(price float64) Trade
}

func (b BuyLimit) price() float64 {
	return b.Bid
}

func (s SellLimit) price() float64 {
	return s.Ask
}

func (b BuyMarket) price() float64 {
	return 1000000.00
}

func (s SellMarket) price() float64 {
	return 0.00
}