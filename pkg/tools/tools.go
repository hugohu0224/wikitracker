package tools

import (
	"sort"
	"sync"
	"wikitracker/pkg/models"
)

func ParseKeyToString(key []byte) string {
	for i, b := range key {
		if b == 0 {
			return string(key[:i])
		}
	}
	return string(key)
}

type CountTracker struct {
	WikiEditInfo map[int64]map[string]*models.WikiEditInfo
	topK         []*models.WikiEditInfo
	Mu           sync.RWMutex
	k            int
}

func NewCountTracker(k int) *CountTracker {
	return &CountTracker{
		WikiEditInfo: make(map[int64]map[string]*models.WikiEditInfo),
		topK:         make([]*models.WikiEditInfo, 0, k),
		Mu:           sync.RWMutex{},
		k:            k,
	}
}

func (ct *CountTracker) GetTopKByStartTime(timestamp int64) []*models.WikiEditInfo {
	ct.Mu.RLock()
	defer ct.Mu.RUnlock()

	topK := make([]*models.WikiEditInfo, ct.k)

	// sub sort
	for _, info := range ct.WikiEditInfo[timestamp] {
		if len(topK) < ct.k {
			topK = append(topK, info)
			if len(topK) == ct.k {
				sort.Slice(topK, func(i, j int) bool {
					return topK[i].EditsCount > topK[j].EditsCount
				})
			}
		} else if info.EditsCount > topK[ct.k-1].EditsCount {
			topK[ct.k-1] = info
			for i := ct.k - 1; i > 0 && topK[i].EditsCount > topK[i-1].EditsCount; i-- {
				topK[i], topK[i-1] = topK[i-1], topK[i]
			}
		}
	}
	return topK
}
