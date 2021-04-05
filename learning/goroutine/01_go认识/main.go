package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func runner(name int) {
	defer wg.Done()
	fmt.Printf("这是第%d子goroutine的运行\n", name)
}

func main() {
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go runner(i)
	}
	wg.Wait()
	fmt.Println("这是主goroutine的运行")
}
