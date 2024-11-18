package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

// 共享内存并发的核心思想是多个 goroutine 共享同一块内存区域（如一个变量、数组或切片），
// 并对其进行读写操作。为了避免数据竞争和保证并发安全，通常需要使用原子操作或者互斥锁来保护共享内存
func Test_one(t *testing.T) {
	var a int             // 共享变量
	var b int32           // 共享变量
	var c int             // 共享变量
	var mu sync.Mutex     // 用于保护共享变量的互斥锁
	var wg sync.WaitGroup // 用于等待所有 Goroutine 完成

	// 启动 1000 个 Goroutine
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// 对 a 进行加 1 操作，不加锁会导致错误
			a++

			atomic.AddInt32(&b, 1)
			mu.Lock()
			c++
			mu.Unlock()

		}()
	}

	// 等待所有 Goroutine 完成
	wg.Wait()

	// 打印结果
	fmt.Println("Final value of a:", a)
	fmt.Println("Final value of b:", b)
	fmt.Println("Final value of c:", c)

}

/*
Final value of a: 96157
Final value of b: 100000
Final value of c: 100000
*/
//管道通信并发,不要共享内存并发，而是通过管道通信的方式并发。
func Test_two(t *testing.T) {
	// 创建一个管道来传递增量
	ch := make(chan int, 10)

	var wg sync.WaitGroup

	// 启动 1000 个 goroutine，将数字 1-5 发送到管道
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ch <- i // 向管道发送i
		}(i)

	}

	// 启动一个 goroutine 接收管道中的数据
	go func() {
		wg.Wait() // 等待所有发送 goroutine 完成

		close(ch) // 完成后关闭管道
	}()
	//range ch 循环会等到管道关闭后才结束,会阻塞主程序直到管道关闭
	for num := range ch {
		fmt.Println("Received:", num)
	}

}
func Test_three(t *testing.T) {
	// 创建两个管道
	ch1 := make(chan int, 50) // 用来发送 1 到 5 的数字
	ch2 := make(chan int, 50) // 用来发送 5 个 1
	var wg sync.WaitGroup

	// 启动 5 个 goroutine，发送数字 1 到 5 到 ch1

	for i := 1; i <= 5; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			ch1 <- i
			ch2 <- i
		}(i)
		go func() {
			defer wg.Done()
			ch2 <- 1
		}()
	}

	// 等待所有 goroutine 完成，sync.WaitGroup 的作用是用来等待多个 goroutine 完成任务。
	go func() {
		wg.Wait()  // 等待所有 goroutine 完成
		close(ch1) // 关闭 ch1
		close(ch2) // 关闭 ch2
		fmt.Println("我的任务完成啦！！")
	}()

	// 使用 select 来同时处理 ch1 和 ch2 的数据
	fmt.Println("开始接收数据：")
	sum := 0
	for {
		select {
		case num, ok := <-ch1:
			if !ok { // 如果 ch1 已经关闭且没有数据了，跳出循环
				ch1 = nil // 关闭 ch1 通道的接收
			} else {
				fmt.Println("Received from ch1:", num)
			}
		case num, ok := <-ch2:
			if !ok { // 如果 ch2 已经关闭且没有数据了，跳出循环
				ch2 = nil // 关闭 ch2 通道的接收
			} else {
				fmt.Println("Received from ch2:", num)
				sum += num // 累加 1 的和
			}
		}
		// 如果 ch1 和 ch2 都关闭了，就退出循环
		if ch1 == nil && ch2 == nil {
			break
		}
	}

	// 打印结果
	fmt.Println("Sum from ch2:", sum)
	fmt.Println("主线程完成啦！！")
}
