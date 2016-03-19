package main

import (
  "fmt"
)

// what makes up a stock market? competing buy and sell orders

type Order struct {
  ticker string
  bid float64
  ask float64
  shares int
  actor string
  filled bool
  canceled bool
}

func main() {
  anOrder := Order{ticker: "GOLANG", bid: 10, shares: 100, actor: "Bob"}
  fmt.Println(anOrder)
}