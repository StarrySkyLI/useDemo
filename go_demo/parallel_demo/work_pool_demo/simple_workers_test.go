package main

import (
	"fmt"
	"testing"
	"time"
)

func TestWorkPool(t *testing.T) {
	task := NewTask(func() error {
		fmt.Print(time.Now())
		return nil
	})
	taskCount := 0
	ticker := time.NewTicker(2 * time.Second)
	p := NewWorkPool(3)
	go func(c *time.Ticker) {
		for {
			p.TaskQueue <- task
			<-c.C
			taskCount++
			if taskCount == 5 {
				p.close()
				break
			}
		}

		return
	}(ticker)
	p.run()
}

type Task struct {
	f func() error
}

func NewTask(f func() error) *Task {
	return &Task{f: f}
}

// Execute 执行业务方法
func (t *Task) Execute() error {
	return t.f()
}

type WorkPool struct {
	TaskQueue chan *Task    // task队列
	workNum   int           // 携程池中最大的worker数量
	shop      chan struct{} // 停止标识
}

// 创建Pool的函数
func NewWorkPool(cap int) *WorkPool {
	if cap <= 0 {
		cap = 10
	}
	return &WorkPool{
		TaskQueue: make(chan *Task),
		workNum:   cap,
		shop:      make(chan struct{}),
	}
}

func (p *WorkPool) worker(workId int) {
	// 具体的工作
	for task := range p.TaskQueue {
		err := task.Execute()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf(" work id %d finished \n", workId)
	}
}

// 携程池开始工作
func (p *WorkPool) run() {
	// 根据work num 去创建worker工作
	for i := 0; i < p.workNum; i++ {
		go p.worker(i)
	}
	<-p.shop
}

func (p *WorkPool) close() {
	p.shop <- struct{}{}
}
