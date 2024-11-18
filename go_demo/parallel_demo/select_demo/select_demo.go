package main

import (
	"fmt"
	"time"
)

var start time.Time

func init() {
	start = time.Now()
}

func service1(c chan string) {
	time.Sleep(3 * time.Second)
	c <- "Hello from service 1"
}

func service2(c chan string) {
	//time.Sleep(5 * time.Second)
	c <- "Hello from service 2"
}

func main() {
	fmt.Println("main start", time.Since(start))

	chan1 := make(chan string)
	chan2 := make(chan string)

	go service1(chan1)
	go service2(chan2)
	//如果所有的 case 语句（通道操作）被阻塞，
	//那么 select 语句将阻塞直到这些 case 条件的一个不阻塞（通道操作），case 块执行。
	//如果有多个 case 块（通道操作）都没有阻塞，
	//那么运行时将随机选择一个不阻塞的 case 块立即执行
	time.Sleep(3 * time.Second)
	select {
	case res := <-chan1:
		fmt.Println("Response form service 1", res, time.Since(start))
	case res := <-chan2:
		fmt.Println("Response form service 2", res, time.Since(start))
		//如果有 case 块的通道操作是非阻塞，那么 select 会执行其 case 块。如果没有那么 select 将默认执行 default 块.
	default:
		fmt.Println("No Response received", time.Since(start))
	}

	fmt.Println("main stop ", time.Since(start))
}

/*
结果：
main start 0s
Response form service 1 Hello from service 1 3.0018445s
main stop  3.0019815s
*/
