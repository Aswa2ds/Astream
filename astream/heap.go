package astream

type asHeap struct {
	data []interface{}
	SortFunc
}

func (heap *asHeap) isEmpty() bool {
	return len(heap.data) == 0
}

func (heap *asHeap) push(elem interface{}) {
	heap.data = append(heap.data, elem)
	heap.up(len(heap.data) - 1)
}

func (heap *asHeap) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !heap.SortFunc(heap.data[i], heap.data[j]) {
			break
		}
		heap.swap(i, j)
		j = i
	}
}

func (heap *asHeap) pop() interface{} {
	n := len(heap.data) - 1
	heap.swap(0, n)
	elem := heap.data[n]
	heap.data = heap.data[:n]
	heap.down(0)
	return elem
}

func (heap *asHeap) swap(i int, j int) {
	heap.data[i], heap.data[j] = heap.data[j], heap.data[i]
}

func (heap *asHeap) down(i0 int) {
	i := i0
	n := len(heap.data)
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && !heap.SortFunc(heap.data[j2], heap.data[j1]) {
			j = j2 // = 2*i + 2  // right child
		}
		if heap.SortFunc(heap.data[j], heap.data[i]) {
			break
		}
		heap.swap(i, j)
		i = j
	}
}
