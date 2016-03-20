package heap

import (
	"testing"
	. "github.com/franela/goblin"
)

func Test(t *testing.T){
	g := Goblin(t)

	g.Describe("Max Heap", func(){

		aHeap := Heap{}

		g.Before(func(){
			aHeap.enqueue(&node{value:1.00, lookup:"bob" })
			aHeap.enqueue(&node{value:2.05, lookup:"bob" })
			aHeap.enqueue(&node{value:2.00, lookup:"bob" })
			aHeap.enqueue(&node{value:1.55, lookup:"bob" })
			aHeap.enqueue(&node{value:0.80, lookup:"bob" })
		})

		g.It("should peek the priority node", func(){
			g.Assert(aHeap.peek().value).Equal(2.05)
		})

		g.It("should dequeue nodes in priority order when dequeing", func(){
			g.Assert(aHeap.dequeue().value).Equal(2.05)
			g.Assert(aHeap.dequeue().value).Equal(2.00)
			g.Assert(aHeap.dequeue().value).Equal(1.55)
			g.Assert(aHeap.dequeue().value).Equal(1.00)
			g.Assert(aHeap.dequeue().value).Equal(0.80)
		})
	})

	g.Describe("Min Heap", func(){
		aHeap := Heap{priority:"min"}

		g.Before(func(){
			aHeap.enqueue(&node{value:1.00, lookup:"bob" })
			aHeap.enqueue(&node{value:2.05, lookup:"bob" })
			aHeap.enqueue(&node{value:2.00, lookup:"bob" })
			aHeap.enqueue(&node{value:1.55, lookup:"bob" })
			aHeap.enqueue(&node{value:0.80, lookup:"bob" })
		})

		g.It("should peek the priority node", func(){
			g.Assert(aHeap.peek().value).Equal(0.80)
		})

		g.It("should dequeue nodes in priority order when dequeing", func(){
			g.Assert(aHeap.dequeue().value).Equal(0.80)
			g.Assert(aHeap.dequeue().value).Equal(1.00)
			g.Assert(aHeap.dequeue().value).Equal(1.55)
			g.Assert(aHeap.dequeue().value).Equal(2.00)
			g.Assert(aHeap.dequeue().value).Equal(2.05)
		})
	})


}