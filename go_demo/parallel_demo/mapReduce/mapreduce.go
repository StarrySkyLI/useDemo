package main

import (
	"fmt"
	"log"

	"github.com/zeromicro/go-zero/core/mr"
)

func demo1() {
	val, err := mr.MapReduce(func(source chan<- int) {
		// generator
		for i := 0; i < 10; i++ {
			source <- i
		}
	}, func(i int, writer mr.Writer[int], cancel func(error)) {
		// mapper
		writer.Write(i * i)
	}, func(pipe <-chan int, writer mr.Writer[int], cancel func(error)) {
		// reducer
		var sum int
		for i := range pipe {
			sum += i
		}
		writer.Write(sum)
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("result:", val)
}
func demo2() {
	generateFunc := func(source chan<- int) {
		for i := 0; i < 10; i++ {
			source <- i
		}
	}

	mapperFunc := func(item int, writer mr.Writer[int], cancel func(error)) {
		writer.Write(item * 2)
	}

	reducerFunc := func(pipe <-chan int, writer mr.Writer[int], cancel func(error)) {
		sum := 0
		for v := range pipe {
			sum += v
		}
		writer.Write(sum)
	}

	result, err := mr.MapReduce(generateFunc, mapperFunc, reducerFunc, mr.WithWorkers(4))
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result) // Output: Result: 90
	}
}
func main() {
	demo1()
	demo2()
}
