package main

import (
	"fmt"
	"strconv"
)

func main() {
	var m int
	fmt.Scan(&m)

	str := strconv.Itoa(m)
	rev := ""

	for i := len(str) - 1; i >= 0; i-- {
		rev += string(str[i])
	}

	ans, _ := strconv.Atoi(rev)
	fmt.Println(ans)
}
