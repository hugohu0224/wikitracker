package initial

import (
	"fmt"
	"time"
	"wikitracker/global"
	"wikitracker/pkg/tools"
)

func InitCountTracker() {
	global.Ct = tools.NewCountTracker(5)
}

func InitCountScan() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("getting topK")
			topK := global.Ct.GetTopKByStartTime(1724296800000)
			fmt.Println(topK)
		}
	}
}

func InitShowTimestamp() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("getting timestamp")
			for key := range global.Ct.WikiEditInfo {
				fmt.Println("timestamp:", key)
			}
		}
	}
}
