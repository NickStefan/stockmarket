package heap

import (

)

type node struct {
	value float64
	lookup string
}

type Heap struct {
	data []*node
}

func (h *Heap) peek() *node {
	if len(h.data) > 0 {
		return h.data[0]
	} else {
		return nil
	}
}

func (h *Heap) insert(node *node) {
 	h.data = append(h.data, node)

	currentNodeIndex := len(h.data) - 1

	if currentNodeIndex == 0 {
		return
	}

	reorder := true

	for reorder {
		parentNodeIndex := h.getParentIndex(currentNodeIndex)

		if h.get(parentNodeIndex).value < h.get(currentNodeIndex).value {
			h.swap(parentNodeIndex, currentNodeIndex)
			currentNodeIndex = parentNodeIndex
		} else {
			reorder = false
		}
	}
}

func (h *Heap) getParentIndex(index int) int {
 	return ((index - 1) / 2);
}

func (h *Heap) get(index int) *node {
	if h.data[index] != nil {
		return h.data[index]
	} else {
		return nil
	}
}

func (h *Heap) swap(indexA int, indexB int){
	a := h.data[indexA]
	b := h.data[indexB]
	h.data[indexA] = b
	h.data[indexB] = a
}
