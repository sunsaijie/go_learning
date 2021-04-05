package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup
var lock sync.Mutex

var x int = 0

func add() {
	defer wg.Done()
	for i := 0; i < 50000; i++ {
		lock.Lock()
		x = x + 1
		lock.Unlock()
	}
}

func main() {
	wg.Add(2)
	go add()
	go add()
	wg.Wait()
	fmt.Println(x)

}
