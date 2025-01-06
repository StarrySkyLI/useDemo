package main

import (
	"container/heap"
	"fmt"
	"sync"
	"time"
)

// 商品数据
type ProductData struct {
	ProductID   string
	AccessCount int
}

// 大顶堆实现
type MaxHeap []ProductData

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i].AccessCount > h[j].AccessCount }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MaxHeap) Push(x interface{}) {
	*h = append(*h, x.(ProductData))
}

func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// 滑动窗口管理
type SlidingWindow struct {
	WindowSize     time.Duration
	Interval       time.Duration
	AccessData     []map[string]int
	Heap           *MaxHeap
	Mutex          sync.Mutex
	CacheThreshold int
	Cache          map[string]int
	LastUpdateTime time.Time
}

// 初始化滑动窗口
func NewSlidingWindow(windowSize, interval time.Duration, cacheThreshold int) *SlidingWindow {
	return &SlidingWindow{
		WindowSize:     windowSize,
		Interval:       interval,
		AccessData:     make([]map[string]int, int(windowSize/interval)),
		Heap:           &MaxHeap{},
		CacheThreshold: cacheThreshold,
		Cache:          make(map[string]int),
		LastUpdateTime: time.Now(),
	}
}

// 添加数据到当前窗口
func (sw *SlidingWindow) AddData(productID string) {
	now := time.Now()
	sw.Mutex.Lock()
	defer sw.Mutex.Unlock()

	// 更新滑动窗口
	if now.Sub(sw.LastUpdateTime) >= sw.Interval {
		sw.LastUpdateTime = now
		sw.rotateWindow()
	}

	// 当前桶索引
	bucketIndex := int(now.Sub(sw.LastUpdateTime) / sw.Interval)
	if sw.AccessData[bucketIndex] == nil {
		sw.AccessData[bucketIndex] = make(map[string]int)
	}

	// 增加该商品的访问次数
	sw.AccessData[bucketIndex][productID]++

	// 重新构建大顶堆
	sw.rebuildHeap()
}

// 滑动窗口滚动，移除过时的数据
func (sw *SlidingWindow) rotateWindow() {
	for i := 1; i < len(sw.AccessData); i++ {
		sw.AccessData[i-1] = sw.AccessData[i]
	}
	sw.AccessData[len(sw.AccessData)-1] = make(map[string]int)
}

// 重新构建大顶堆
func (sw *SlidingWindow) rebuildHeap() {
	sw.Heap = &MaxHeap{}
	for _, bucket := range sw.AccessData {
		for productID, count := range bucket {
			heap.Push(sw.Heap, ProductData{ProductID: productID, AccessCount: count})
		}
	}
}

// 获取TopK热点数据并升级到缓存
func (sw *SlidingWindow) UpdateCache(k int) {
	sw.Mutex.Lock()
	defer sw.Mutex.Unlock()

	topK := []ProductData{}
	for i := 0; i < k && sw.Heap.Len() > 0; i++ {
		topProduct := heap.Pop(sw.Heap).(ProductData)
		topK = append(topK, topProduct)

		// 判断是否超过访问阈值，加入缓存
		if topProduct.AccessCount >= sw.CacheThreshold {
			sw.Cache[topProduct.ProductID] = topProduct.AccessCount
			fmt.Printf("Product %s has been cached with %d accesses.\n", topProduct.ProductID, topProduct.AccessCount)
		}
	}

	// 输出TopK商品
	fmt.Println("Top K products:")
	for _, product := range topK {
		fmt.Printf("ProductID: %s, AccessCount: %d\n", product.ProductID, product.AccessCount)
	}
}

func main() {
	// 定义滑动窗口为10秒，滑动间隔为2秒，缓存阈值为500次
	sw := NewSlidingWindow(10*time.Second, 2*time.Second, 500)

	// 模拟访问商品
	productIDs := []string{"product1", "product2", "product3", "product4", "product5"}
	for i := 0; i < 1000; i++ {
		productID := productIDs[i%len(productIDs)]
		sw.AddData(productID)
		time.Sleep(10 * time.Millisecond)
	}

	// 更新并打印热点商品，获取Top3
	sw.UpdateCache(3)
}
