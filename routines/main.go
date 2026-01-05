package main

import (
	"fmt"
	"time"
)

func main() {
	go greet("hello")
	go greet("world")
	time.Sleep(10 * time.Second)
}

func greet(s string) {
	for i := 0; i < 5; i++ {
		fmt.Println(s)
	}
}
