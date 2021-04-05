package main

import (
	"fmt"
	"sync"
	"time"
)

var x int = 0
var rwlock sync.RWMutex
var wg sync.WaitGroup

func write() {
	defer wg.Done()
	rwlock.Lock()
	x = x + 1
	time.Sleep(time.Microsecond * 10)
	rwlock.Unlock()
}

func read() {
	defer wg.Done()
	rwlock.RLock()
	a := x
	fmt.Println(a)
	time.Sleep(time.Microsecond * 5)
	rwlock.RUnlock()

}

func main() {
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go write()
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go read()
	}
	wg.Wait()

}
