package pkg

type Heap struct {
	Objects []interface{}
	Less    func(a, b interface{}) bool
}

func NewHeap(less func(a, b interface{}) bool) *Heap {
	return &Heap{
		Objects: make([]interface{}, 1),
		Less:    less,
	}
}

func (h *Heap) Push(obj interface{}) {
	h.Objects = append(h.Objects, obj)
	for i := len(h.Objects) - 1; i > 1 && h.Less(h.Objects[i], h.Objects[i/2]); i /= 2 {
		h.Objects[i], h.Objects[i/2] = h.Objects[i/2], h.Objects[i]
	}
}

func (h *Heap) Pop() interface{} {
	if len(h.Objects) == 1 {
		return nil
	}
	return h.down(1)
}

func (h *Heap) Peek() interface{} {
	if len(h.Objects) == 1 {
		return nil
	}
	return h.Objects[1]
}

func (h *Heap) Length() int {
	return len(h.Objects) - 1
}

func (h *Heap) down(i int) (obj interface{}) {
	obj, h.Objects[i] = h.Objects[i], h.Objects[len(h.Objects)-1]
	h.Objects = h.Objects[:len(h.Objects)-1]

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
	return
}
