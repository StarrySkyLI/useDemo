package main

import (
	"fmt"
	"testing"
)

func square(c chan int) {
	fmt.Println("[square] reading")
	num := <-c
	c <- num * num
}

func cube(c chan int) {
	fmt.Println("[cube] reading")
	num := <-c
	c <- num * num * num
}

// 写两个协程，一个用来计算数字的平方，另一个用来计算数字的立方。
func Test_two_go(t *testing.T) {
	fmt.Println("[main] main started")
	squareChan := make(chan int)
	cubeChan := make(chan int)
	go square(squareChan)
	go cube(cubeChan)

	testNum := 3

	fmt.Println("[main] send testNum to squareChan")
	squareChan <- testNum
	fmt.Println("[main] resuming")

	fmt.Println("[main] send testNum to cubeChane")
	cubeChan <- testNum
	fmt.Println("[main] resuming")

	fmt.Println("[main] reading from channels")
	squareVal, cubeVal := <-squareChan, <-cubeChan
	sum := squareVal + cubeVal

	fmt.Println("[main] sum of square and cube of", testNum, " is", sum)
	fmt.Println("[main] main stop")
}

/*
[main] main started
[main] send testNum to squareChan
[cube] reading
[square] reading
[main] resuming
[main] send testNum to cubeChane
[main] resuming
[main] reading from channels
[main] sum of square and cube of 3  is 36
[main] main stop
*/
