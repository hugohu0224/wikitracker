package tools

import (
	"container/heap"
	"strconv"
)

func ParseKeyToString(key []byte) string {
	for i, b := range key {
		if b == 0 {
			return string(key[:i])
		}
	}
	return string(key)
}

type Item struct {
	title string
	count int
	url   string
}

type MaxHeap []Item

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i].count > h[j].count }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) {
	*h = append(*h, x.(Item))
}

func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func TopK(editCounts map[string]int, k int) []map[string]string {
	h := &MaxHeap{}
	heap.Init(h)

	for title, count := range editCounts {
		heap.Push(h, Item{title: title, count: count})
	}

	result := make([]map[string]string, 0, k)
	for i := 0; i < k && h.Len() > 0; i++ {
		item := heap.Pop(h).(Item)
		resultMap := map[string]string{
			"title": item.title,
			"count": strconv.Itoa(item.count),
		}
		result = append(result, resultMap)
	}
	return result
}
