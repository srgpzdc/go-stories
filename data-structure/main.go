package main

import (
	"data-structure/heap"
	"fmt"
)

func main() {
	myHeap := heap.New(1.2, 2.1, 3.1, -100.1)
	size := len(*myHeap)
	fmt.Printf("Heap size: %d\n", size)
	fmt.Printf("%v\n", myHeap)

	myHeap.Push(float32(-100.2))
	myHeap.Push(float32(0.2))
	fmt.Printf("Heap size: %d\n", len(*myHeap))
	fmt.Printf("%v\n", myHeap)
	fmt.Printf("%v\n", myHeap)

	for myHeap.Len() > 0 {
		fmt.Println(myHeap.Pop())
	}
}
