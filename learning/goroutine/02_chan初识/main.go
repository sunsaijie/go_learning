package main

import (
	"fmt"
	"math/rand"
)

// 随机生成数字，放入task通道,
// 新建一个20个goroutine的池，计算该数字 * 2，将结果存入result通道
// 主goroutine从result通道中取值后，打印结果

var task chan int64
var result chan int64

func calc(worker int, task <-chan int64, result chan<- int64) {
	for {
		v, ok := <-task
		if !ok {
			break
		}
		fmt.Printf("第%d个worker开始工作\n", worker)
		ret := v * 2
		result <- ret
		fmt.Printf("第%d个worker结束工作\n", worker)
	}
}

func main() {
	task = make(chan int64, 100)
	result = make(chan int64, 100)
	for t := 0; t < 5; t++ {
		go calc(t, task, result)
	}
	for i := 0; i < 100; i++ {
		a := rand.Int63()
		task <- int64(a)
	}
	close(task)

	for {
		a := <-result
		fmt.Printf("结果为%d\n", a)
	}

}
