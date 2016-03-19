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
		aHeap.insert(&node{value:2.05, lookup:"bob" })
		aHeap.insert(&node{value:2.00, lookup:"bob" })
		aHeap.insert(&node{value:1.55, lookup:"bob" })
		aHeap.insert(&node{value:0.80, lookup:"bob" })

		g.Assert(aHeap.peek().value).Equal(2.05)
	})

	g.It("should dequeue nodes in priority order when dequeing", func(){

	})

  })
}