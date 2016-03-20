package heap

import (
	"testing"
	. "github.com/franela/goblin"
)

func Test(t *testing.T){
	g := Goblin(t)

	g.Describe("Max Heap", func(){

		aHeap := Heap{Priority:"Max"}

		g.Before(func(){
			aHeap.Enqueue(&Node{Value:1.00, Lookup:"bob" })
			aHeap.Enqueue(&Node{Value:2.05, Lookup:"bob" })
			aHeap.Enqueue(&Node{Value:2.00, Lookup:"bob" })
			aHeap.Enqueue(&Node{Value:1.55, Lookup:"bob" })
			aHeap.Enqueue(&Node{Value:0.80, Lookup:"bob" })
		})

		g.It("should peek the priority Node", func(){
			g.Assert(aHeap.Peek().Value).Equal(2.05)
		})

		g.It("should dequeue Nodes in priority order when dequeing", func(){
			g.Assert(aHeap.Dequeue().Value).Equal(2.05)
			g.Assert(aHeap.Dequeue().Value).Equal(2.00)
			g.Assert(aHeap.Dequeue().Value).Equal(1.55)
			g.Assert(aHeap.Dequeue().Value).Equal(1.00)
			g.Assert(aHeap.Dequeue().Value).Equal(0.80)
		})
	})

	g.Describe("Min Heap", func(){
		aHeap := Heap{Priority:"min"}

		g.Before(func(){
			aHeap.Enqueue(&Node{Value:1.00, Lookup:"bob" })
			aHeap.Enqueue(&Node{Value:2.05, Lookup:"bob" })
			aHeap.Enqueue(&Node{Value:2.00, Lookup:"bob" })
			aHeap.Enqueue(&Node{Value:1.55, Lookup:"bob" })
			aHeap.Enqueue(&Node{Value:0.80, Lookup:"bob" })
		})

		g.It("should peek the priority Node", func(){
			g.Assert(aHeap.Peek().Value).Equal(0.80)
		})

		g.It("should dequeue Nodes in priority order when dequeing", func(){
			g.Assert(aHeap.Dequeue().Value).Equal(0.80)
			g.Assert(aHeap.Dequeue().Value).Equal(1.00)
			g.Assert(aHeap.Dequeue().Value).Equal(1.55)
			g.Assert(aHeap.Dequeue().Value).Equal(2.00)
			g.Assert(aHeap.Dequeue().Value).Equal(2.05)
		})
	})


}