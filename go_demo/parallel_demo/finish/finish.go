package main

import (
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/timex"
	"time"

	"github.com/zeromicro/go-zero/core/mr"
)

func Demo1() {
	start := timex.Now()

	mr.FinishVoid(func() {
		time.Sleep(time.Second)
	}, func() {
		time.Sleep(time.Second * 5)
	}, func() {
		time.Sleep(time.Second * 10)
	}, func() {
		time.Sleep(time.Second * 6)
	}, func() {
		if err := mr.Finish(func() error {
			time.Sleep(time.Second)
			return nil
		}, func() error {
			time.Sleep(time.Second * 10)
			return nil
		}); err != nil {
			fmt.Println(err)
		}
	})

	fmt.Println(timex.Since(start))
}
func demo2() {
	funcs := []func() error{
		func() error {
			fmt.Println("Function 1 executed")
			return nil
		},
		func() error {
			fmt.Println("Function 2 executed")
			return errors.New("error in function 2")
		},
	}

	err := mr.Finish(funcs...)
	if err != nil {
		fmt.Println("Finish encountered an error:", err)
	}

	voidFuncs := []func(){
		func() {
			fmt.Println("Void Function 1 executed")
		},
		func() {
			fmt.Println("Void Function 2 executed")
		},
	}

	mr.FinishVoid(voidFuncs...)
}
func main() {
	Demo1()
	demo2()
}
