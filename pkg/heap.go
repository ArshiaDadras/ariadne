package pkg

type Heap struct {
	Objects map[int]interface{}
	Less    func(a, b interface{}) bool
}

func NewHeap(less func(a, b interface{}) bool) *Heap {
	return &Heap{
		Objects: make(map[int]interface{}),
		Less:    less,
	}
}

func (h *Heap) Push(obj interface{}) {
	h.Objects[len(h.Objects)+1] = obj
	h.up(len(h.Objects))
}

func (h *Heap) Pop() interface{} {
	if len(h.Objects) == 0 {
		return nil
	}

	obj := h.Objects[1]
	h.Objects[1] = h.Objects[len(h.Objects)]
	delete(h.Objects, len(h.Objects))
	h.down(1)

	return obj
}

func (h *Heap) Len() int {
	return len(h.Objects)
}

func (h *Heap) up(i int) {
	for i > 1 && h.Less(h.Objects[i], h.Objects[i/2]) {
		h.Objects[i], h.Objects[i/2] = h.Objects[i/2], h.Objects[i]
		i /= 2
	}
}

func (h *Heap) down(i int) {
	for 2*i < len(h.Objects) {
		j := 2 * i
		if j < len(h.Objects)-1 && h.Less(h.Objects[j+1], h.Objects[j]) {
			j++
		}

		if !h.Less(h.Objects[j], h.Objects[i]) {
			break
		}

		h.Objects[i], h.Objects[j] = h.Objects[j], h.Objects[i]
		i = j
	}
}
