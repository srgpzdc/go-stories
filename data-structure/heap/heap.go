package heap

import "container/heap"

type heapFloat32 []float32

func (n *heapFloat32) Pop() interface{} {
	old := *n
	x := old[len(old)-1]
	new := old[0 : len(old)-1]
	*n = new
	return x
}

func (n *heapFloat32) Push(x interface{}) {
	*n = append(*n, x.(float32))
	heap.Init(n)
}

func (n heapFloat32) Len() int {
	return len(n)
}
func (n heapFloat32) Less(a, b int) bool {
	return n[a] < n[b]
}
func (n heapFloat32) Swap(a, b int) {
	n[a], n[b] = n[b], n[a]
}

func New(items ...float32) *heapFloat32 {
	var myHeap heapFloat32 = items
	heap.Init(&myHeap)
	return &myHeap
}
