package main

import (
  "fmt"
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

type OrderBook {

}

// whats an algorithm to match buyers with sellers? simple case just using market orders

func main() {
  anOrder := Order{ticker: "GOLANG", bid: 10, shares: 100, actor: "Bob"}
  fmt.Println(anOrder)
}