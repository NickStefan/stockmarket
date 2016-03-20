package heap

type node struct {
	value float64
	lookup string
}

type Heap struct {
	priority string
	data []*node
}

func (h *Heap) peek() *node {
	if len(h.data) > 0 {
		return h.data[0]
	} else {
		return nil
	}
}

func (h *Heap) compare(current *node, other *node) bool {
	if h.priority == "min" {
		// in a min heap,
			// keep swapping when this is true during the enqueue
			// and stop swapping when this is true in the denqueue reorder
		return other.value > current.value
	} else {
		// in a max heap,
			// keep swapping when this is true during the enqueue
			// and stop swapping when this is true in the denqueue reorder
		return other.value < current.value
	}
}

func (h *Heap) enqueue(node *node) {
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

func (h *Heap) dequeue () *node {
	if !(len(h.data) > 0) {
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

func (h *Heap) get(index int) *node {
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
