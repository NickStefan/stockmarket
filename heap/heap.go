package heap

import (
	"fmt"
)

type Node struct {
	Value float64
	Lookup string
}

type Heap struct {
	Priority string
	data []*Node
}

func (h *Heap) Peek() *Node {
	if len(h.data) > 0 {
		return h.data[0]
	} else {
		return nil
	}
}

func (h *Heap) compare(current *Node, other *Node) bool {
	if h.Priority == "min" {
		// in a min heap,
			// keep swapping when this is true during the enqueue
			// and stop swapping when this is true in the denqueue reorder
		return other.Value > current.Value
	} else {
		// in a max heap,
			// keep swapping when this is true during the enqueue
			// and stop swapping when this is true in the denqueue reorder
		return other.Value < current.Value
	}
}

func (h *Heap) Enqueue(node *Node) {
 	h.data = append(h.data, node)

	currentNode := len(h.data) - 1

	if currentNode == 0 {
		return
	}

	reorder := true

	for reorder {
		parentNode := h.getParent(currentNode)

		if h.compare(h.get(currentNode), h.get(parentNode)) {
			h.swap(parentNode, currentNode)
			currentNode = parentNode
		} else {
			reorder = false
		}
	}
}

func (h *Heap) Dequeue () *Node {
	if !(len(h.data) > 0) {
		fmt.Println("null!")
		return nil
	}
	dequeued := h.data[0]

	if len(h.data) > 0 {
		h.swap(0, len(h.data) - 1)
		h.data = h.data[0:(len(h.data) - 1)]
	} else {
		return nil
	}

	current := 0
	reorder := true

	for reorder {
		leftChild, rightChild := h.getChildren(current)
		
		leftChildNode := h.get(leftChild)
		rightChildNode := h.get(rightChild)

		var minChild int

		if leftChildNode == nil && rightChildNode != nil {
			minChild = rightChild
		
		} else if rightChildNode == nil && leftChildNode != nil {
			minChild = leftChild
		
		} else if rightChildNode != nil && leftChildNode != nil {
			if h.compare(leftChildNode, rightChildNode) {
				minChild = leftChild
			} else {
				minChild = rightChild
			}
		}

		if minChild == 0 || h.compare(h.get(current), h.get(minChild)) {
			reorder = false
		} else	{
			h.swap(minChild, current)
			current = minChild
		}
	}

	return dequeued
}

func (h *Heap) getParent(index int) int {
 	return ((index - 1) / 2);
}

func (h *Heap) getChildren(index int) (int, int) {
	return (2 * index + 1), (2 * index + 2)
}

func (h *Heap) get(index int) *Node {
	if index >= len(h.data) {
		return nil
	} else {
		return h.data[index]
	}
}

func (h *Heap) swap(indexA int, indexB int){
	a := h.data[indexA]
	b := h.data[indexB]
	h.data[indexA] = b
	h.data[indexB] = a
}
