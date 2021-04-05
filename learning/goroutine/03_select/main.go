package main

import "fmt"

func main() {
	ch1 := make(chan int, 1)
	for i := 0; i < 10; i++ {
		select {
		case ch1 <- i:
			fmt.Println("存入", i)
		case a := <-ch1:
			fmt.Printf("取值%d\n", a)
		}
	}
}
