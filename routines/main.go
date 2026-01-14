package main

import (
	"fmt"
	"net/http"
	"routines/advance"
	"sync"
)

var signal = []string{"test"}
var wg sync.WaitGroup
var mut sync.Mutex

func main() {
	// endpoints := []string{
	// 	"https://youtube.com",
	// 	"https://google.com",
	// 	"https://github.com",
	// 	"https://go.dev",
	// 	"https://fb.com",
	// }
	// for _, web := range endpoints {
	// 	go getStatusCode(web)
	// 	wg.Add(1)
	// }
	// wg.Wait()
	// fmt.Println(signal)
	// advance.Mutext()
	advance.Channels()
}

//	func greet(s string) {
//		for i := 0; i < 5; i++ {
//			fmt.Println(s)
//		}
//	}
func getStatusCode(endpoint string) {
	defer wg.Done()
	res, err := http.Get(endpoint)
	if err != nil {
		fmt.Println("oops in endpoint")
	}
	mut.Lock()
	signal = append(signal, endpoint)
	mut.Unlock()
	fmt.Printf("%d status code for %s", res.StatusCode, endpoint)
	fmt.Println()

}
