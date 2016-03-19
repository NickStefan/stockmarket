package heap

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
}
