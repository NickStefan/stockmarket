package heap

import (
  "testing"
  . "github.com/franela/goblin"
)

func Test(t *testing.T){
  g := Goblin(t)

  g.Describe("Heap", func(){

    g.It("should peek the priority node", func(){
      aHeap := Heap{}
      aHeap.insert(&node{value:1.00, lookup:"bob" })
      g.Assert(aHeap.peek().value).Equal(1.00)
    })

    g.It("should dequeue nodes in priority order when dequeing", func(){

    })

  })
}