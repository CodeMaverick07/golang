package advance

import (
	"fmt"
	"sync"
)

func Channels() {
	fmt.Println("channels")
	myCh := make(chan int, 2)
	wg := &sync.WaitGroup{}
	wg.Add(2)
	// <-chan recive only chan
	go func(ch <-chan int, wg *sync.WaitGroup) {
		val, isChanelOpen := <-myCh
		fmt.Println(isChanelOpen)
		fmt.Println(val)
		wg.Done()
	}(myCh, wg)
	// chan<- send only chan
	go func(ch chan<- int, wg *sync.WaitGroup) {

		myCh <- 5
		myCh <- 6
		close(myCh)
		wg.Done()
	}(myCh, wg)

	wg.Wait()
}
