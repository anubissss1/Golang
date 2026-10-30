package main

import (
	"fmt"
	"os"
	"strconv"
)

//go:noinline
func depth(n int) int {
	var pad [64]byte // Change to [1024]byte for task 1, part 1.
	pad[n%len(pad)] = byte(n)
	if n == 0 {
		return 0
	}
	return depth(n-1) + 1 + int(pad[(n+1)%len(pad)])
}

//go:noinline
func forever(n int) int { return forever(n+1) + 1 }

func main() {
	if len(os.Args) > 1 && os.Args[1] == "forever" {
		fmt.Println(forever(0))
		return
	}
	n := 1000000
	if len(os.Args) > 1 {
		v, err := strconv.Atoi(os.Args[1])
		if err == nil && v >= 0 {
			n = v
		}
	}
	fmt.Println("depth:", depth(n))
}
