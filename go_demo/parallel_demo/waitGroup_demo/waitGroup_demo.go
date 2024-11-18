package main

import (
	"fmt"
	"sync"
	"time"
)

func service(wg *sync.WaitGroup, instance int) {
	time.Sleep(2 * time.Second)
	fmt.Println("Service called on instance", instance)

	wg.Done() //协程数-1
}

func main() {
	fmt.Println("main started")
	//有一种业务场景是你需要知道所有的协程是否已执行完成他们的任务。
	//这个和只需要随机选择一个条件为true 的 select 不同，
	//他需要你满足所有的条件都是 true 才可以激活主线程继续执行。 这里的条件指的是非阻塞的通道操作。

	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {

		wg.Add(1)

		go service(&wg, i)
	}
	//在 for 循环执行完成后，我们通过调用 wg.Wait() 去阻塞当前主线程，并把调度权让给其他协程，
	//直到计数器值为 0 之后，主线程才会被再次调度。
	wg.Wait() //阻塞
	fmt.Println("main stop ")
}

/*
结果：(结果是不唯一的，一共有3!次可能的结果)
main started
Service called on instance 2
Service called on instance 1
Service called on instance 3
main stop
*/
