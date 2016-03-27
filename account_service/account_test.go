package main

import (
	. "github.com/franela/goblin"
	"testing"
)

func test(t *testing.T){
	g := Goblin(t)

	g.Describe("Account", func(){

		g.It("should increment and decrement shares in Assets hash", func(){

		})

		g.It("should update the cash balance", func(){

		})

		// g.It("should recalculate the account value on every trade")
	})
}