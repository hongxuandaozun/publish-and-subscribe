package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"hongxuan.com/pubsub/pubsub"
)

var total struct {
	sync.Mutex
	value int
}

func worker(wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i <= 100; i++ {
		total.Lock()
		total.value += i
		total.Unlock()
	}
}

var work = []func(){
	func() { println("1"); time.Sleep(1 * time.Second) },
	func() { println("2"); time.Sleep(1 * time.Second) },
	func() { println("3"); time.Sleep(1 * time.Second) },
	func() { println("4"); time.Sleep(1 * time.Second) },
	func() { println("5"); time.Sleep(1 * time.Second) },
}
// 生产者: 生成 factor 整数倍的序列
func Producer(factor int, out chan<- int) {
	for i := 0; ; i++ {
		out <- i*factor
	}
}

// 消费者
func Consumer(in <-chan int) {
	for v := range in {
		fmt.Println(v)
	}
}
func main() {

	p := pubsub.NewPublisher(time.Millisecond*100, 10)
	defer p.Close()
	all := p.Subscribe()
	golang := p.SubscribeTopic(func(v interface{}) bool {
		if s,ok := v.(string);ok{
			return strings.Contains(s,"golang")
		}
		return false
	})
	p.Publish("hello,  world!")
	p.Publish("hello, golang!")

	go func() {
		for  msg := range all {
			fmt.Println("all:", msg)
		}
	} ()

	go func() {
		for  msg := range golang {
			fmt.Println("golang:", msg)
		}
	} ()

	// 运行一定时间后退出
	time.Sleep(3 * time.Second)
	//ch := make(chan int, 64) // 成果队列
	//
	//go Producer(3, ch) // 生成 3 的倍数的序列
	//go Producer(5, ch) // 生成 5 的倍数的序列
	//go Consumer(ch)    // 消费 生成的队列
	//
	//// Ctrl+C 退出
	//sig := make(chan os.Signal, 1)
	//signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	//fmt.Printf("quit (%v)\n", <-sig)
}