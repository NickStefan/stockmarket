package main

import (
 	"strconv"
)

// what makes up a stock market? competing buy and sell orders

// market buy order "i will buy 100 shares at the lowest available sell price"
// limit buy order "i will buy 100 shares at the lowest available price below ____"

// market sell order "i will sell 100 shares at the highest available price"
// limit sell order "i will sell 100 shares highest available price above _____"

type BaseOrder struct {
	shares int
	ticker string
	actor string // BOB
	intent string // BUY || SELL
	kind string // MARKET || LIMIT
	state string // OPEN || FILLED || CANCELED
	timecreated int64 // unix time
}

func (b *BaseOrder) lookup() string {
	return b.actor + strconv.FormatInt(b.timecreated, 10)
}

func (b *BaseOrder) getOrder() *BaseOrder {
	return b
}

func (b *BaseOrder) partialFill(price float64, newShares int) Trade {
	b.shares = newShares
	return Trade{
		Actor: b.actor, Shares: b.shares - newShares,
		Price: price, Intent: b.intent, Kind: b.kind, Ticker: b.ticker,
	}
}

func (b *BaseOrder) fill(price float64) Trade {
	return Trade{
		Actor: b.actor, Shares: b.shares, Price: price,
		Intent: b.intent, Ticker: b.ticker, Kind: b.kind,
	}
}

type BuyLimit struct {
	bid float64
	*BaseOrder
}

type SellLimit struct {
	ask float64
	*BaseOrder
}

type BuyMarket struct {
	*BaseOrder
}

type SellMarket struct {
	*BaseOrder
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
	return b.bid
}

func (s SellLimit) price() float64 {
	return s.ask
}

func (b BuyMarket) price() float64 {
	return 1000000.00
}

func (s SellMarket) price() float64 {
	return 0.00
}