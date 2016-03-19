package main

import (
  "fmt"
  "github.com/nickstefan/market/heap"
)

// what makes up a stock market? competing buy and sell orders

// market buy order "i will buy 100 shares at the lowest available sell price"
// limit buy order "i will buy 100 shares at the lowest available price below ____"

// market sell order "i will sell 100 shares at the highest available price"
// limit sell order "i will sell 100 shares highest available price above _____"

type Order struct {
  ticker string
  bid float64
  ask float64
  shares int
  actor string
  filled bool
  canceled bool
}

// how does a stock market organize the orders? Depth of Market or OrderBook

type OrderBook struct {

}

// whats an algorithm to match buyers with sellers? simple case just using market orders

// orderBook.add( *order ) 
// puts the order into the buy or sell sub data structures
// then the order book attempts to fill that order
// if can fulfill
//    fills order, and removes the buy and sell orders that were affected by the fill
// 
// repeat on every orderBook.add( *order )

// orderBook.remove( *order )
// removes order from proper sub data structure

// so what does it mean to fill an order, and how could two data structures help us do that?

// an overlap between two heaps prioritized by price?
// but do we need constant time access to things deep inside the heap,
// due to order remove()?
// no, we mainly just need to pull things from the top of the buy heap and sell heap

// filling an order could be as simple as peeking at the buy heap and sell heap
// and asking if buy heap limit >= sell heap limit

// what about market orders? how would they prioritize into the heaps?
// do we need 4 heaps? or would market orders just always have priority?

// I think two heaps with market having highest priority works

// how would a limit order ever get filled if market is always higher priority?
// say the highest buy order is for 1000 shares limit $10.00,
// say the highest sell order is a market for 100 shares
// say the next highest sell order is a limit for 600 at $9.98.

// the trade would get partially filled at 100 shares @ $10.00 and
// 600 shares @ 9.99

func main() {
  anOrder := Order{ticker: "GOLANG", bid: 10, shares: 100, actor: "Bob"}
  aHeap := heap.Heap{}
  fmt.Println(anOrder)
  fmt.Println(aHeap)
}