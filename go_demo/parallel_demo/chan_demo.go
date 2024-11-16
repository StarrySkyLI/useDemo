package main

import (
	"fmt"
	"time"
)

// 声明一个函数 greet, 这个函数的参数 c 是一个 string 类型的通道。
// 在这个函数中，我们从通道 c 中接收数据并打印到控制台上。
func greet(c chan string) {
	//for这块正常执行
	for i := 3; i > 0; i-- {
		time.Sleep(1 * time.Second)
		fmt.Println("greet倒计时 :", i, "s")
	}
	fmt.Println("greet等待读channel数据阻塞")
	//如果当前协程正在从一个没有任何值的通道中读取数据，那么当前协程会阻塞并且等待其他协程往此通道写入值。
	//因此，读操作将被阻塞
	for val := range c {
		fmt.Println("Hello " + val + "!")
	}

	/*
		注释该行会导致没有协程接收管道的值
		报错
		main start
		fatal error: all goroutines are asleep - deadlock!  //所有协程都进入休眠状态,死锁

		goroutine 1 [chan send]:
		main.main()
	*/

}

func main() {
	fmt.Println("Main Start")
	// main 函数的第一个语句是打印 main start 到控制台。
	channel := make(chan string)
	// 在 main 函数中使用 make 函数创建一个 string 类型的通道赋值给 ‘ channel ’ 变量
	go greet(channel)

	fmt.Println("这时候greet协程读channel那部分是阻塞的，直到有值传入")
	for i := 10; i > 0; i-- {
		time.Sleep(1 * time.Second)
		fmt.Println("main倒计时 :", i, "s")
	}
	// 把 channel 通道传递给 greet 函数并用 go 关键词以协程方式运行它。
	// 此时，程序有两个协程并且正在调度运行的是 main goroutine 主函数
	// 给通道 channel 传入一个数据 DEMO.
	// 此时主线程将阻塞直到有协程接收这个数据. Go 的调度器开始调度 greet 协程接收通道 channel 的数据
	fmt.Println("主线程阻塞，等待greet读数据")
	channel <- "DEMO1"
	time.Sleep(1 * time.Second)
	fmt.Println("主线程阻塞，等待greet读数据")
	channel <- "DEMO2"
	time.Sleep(1 * time.Second)
	close(channel)
	//不能向一个关了的channel发信息
	//channel <- "DEMO3"

	fmt.Println("Main Stop")
	// 然后主线程激活并且执行后面的语句，打印 main stopped
}

/*
Main Start
这时候greet协程读channel那部分是阻塞的，直到有值传入
main倒计时 : 10 s
greet倒计时 : 3 s
greet倒计时 : 2 s
main倒计时 : 9 s
main倒计时 : 8 s
greet倒计时 : 1 s
greet等待读channel数据阻塞
main倒计时 : 7 s
main倒计时 : 6 s
main倒计时 : 5 s
main倒计时 : 4 s
main倒计时 : 3 s
main倒计时 : 2 s
main倒计时 : 1 s
主线程阻塞，等待greet读数据
Hello DEMO1!
主线程阻塞，等待greet读数据
Hello DEMO2!
Main Stop


*/
