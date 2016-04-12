package main

import (
  "testing"
  . "github.com/franela/goblin"
)

func Test(t *testing.T) {
  g := Goblin(t)

  g.Describe("Ticker", func(){

    g.It("QUOTE channel, should receive realtime trades", func(){

    })

    g.It("CHART channel, should receive realtime ticks", func(){

    })

    g.It("on CHART, upstream message nth history, should publish nth history of ticks", func(){

    })
  }) 
}