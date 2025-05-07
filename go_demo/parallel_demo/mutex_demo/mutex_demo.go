package main

import (
	"fmt"
	"sync"
	"time"
)

var mu sync.Mutex
var mu2 sync.Mutex
var counter int
var batch = 1

func main() {
	counter = 0
	var wg sync.WaitGroup

	// 启动多个 goroutine 来模拟并发访问
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			s := service(1)
			fmt.Println(s)
		}(i)
	}

	// 等待所有 goroutine 完成
	wg.Wait()

}
func service(id int) int {
	// time.Sleep(500 * time.Millisecond)
	//
	// if counter > 0 {
	// 	counter--
	// }
	//
	// if counter == 0 {
	// 	add()
	// }
	// fmt.Printf("id:%d batch :%d  counter: %d\n ", id, batch, counter)
	if counter == 0 {
		if mu.TryLock() {
			defer mu.Unlock()
			counter++
		}
	}

	return counter

}
func add() {

	time.Sleep(1 * time.Second)
	batch++
	counter = 20
	fmt.Println("以生成下一批：", batch)
}
