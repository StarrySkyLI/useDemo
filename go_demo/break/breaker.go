package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/zeromicro/go-zero/core/breaker"
)

type mockError struct {
	status int
}

func (e mockError) Error() string {
	return fmt.Sprintf("HTTP STATUS: %d", e.status)
}

func main() {
	for i := 0; i < 100; i++ {

		if err := breaker.DoWithFallbackAcceptable("test", func() error {
			fmt.Println("当前次数：", i)
			return mockRequest()
		}, func(err error) error {
			fmt.Println("在：", i, "，次发生熔断")
			//发生了熔断，这里可以自定义熔断错误转换
			return errors.New("当前服务不可用，请稍后再试")
		}, func(err error) bool { // 当 mock 的http 状态码不为500时都会被认为是正常的，否则加入错误窗口
			me, ok := err.(mockError)
			if ok {
				return me.status != 500
			}
			return false
		}); err != nil {
			println(err.Error())
		}
	}

}

func mockRequest() error {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	num := r.Intn(100)

	if num%4 == 0 {
		fmt.Println("rand num 正常:", num)
		return nil
	} else if num%5 == 0 {
		fmt.Println("rand num 500错误:", num)
		return mockError{status: 500}
	}
	fmt.Println("rand num 普通错误:", num)
	return errors.New("dummy")
}
