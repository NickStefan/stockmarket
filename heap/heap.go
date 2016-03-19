package heap

type node struct {
  value float64
  lookup string
  leftChild *node
  rightChild *node
}

type Heap struct {
  priorityNode *node
}

func (h *Heap) peek() *node {
  return h.priorityNode
}

func (h *Heap) insert(node *node) {
  if h.priorityNode == nil {
    h.priorityNode = node
  }
}

func (h *Heap) dequeue() *node {
  return h.priorityNode
}
