package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// Each row of the handout's 3x4 channel table can be run independently.
func cell(state, op string) {
	var ch chan int
	switch state {
	case "open":
		ch = make(chan int, 1)
	case "closed":
		ch = make(chan int, 1)
		close(ch)
	case "nil": // deliberately nil
	default:
		panic("unknown state")
	}
	switch op {
	case "send":
		ch <- 7
		fmt.Println("send completed")
	case "receive":
		v, ok := <-ch
		fmt.Println("received:", v, "ok:", ok)
	case "close":
		close(ch)
		fmt.Println("close completed")
	case "size":
		fmt.Printf("len=%d cap=%d\n", len(ch), cap(ch))
	}
}
func ping(receiverSleeps bool) {
	start := time.Now()
	ch := make(chan string)
	done := make(chan struct{})
	go func() {
		defer close(done)
		if !receiverSleeps {
			time.Sleep(500 * time.Millisecond)
		}
		fmt.Println("sender: about to send at", time.Since(start).Round(10*time.Millisecond))
		ch <- "ping"
		fmt.Println("sender: send finished at", time.Since(start).Round(10*time.Millisecond))
	}()
	if receiverSleeps {
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println("receiver: waiting at", time.Since(start).Round(10*time.Millisecond))
	fmt.Println("receiver: got", <-ch, "at", time.Since(start).Round(10*time.Millisecond))
	<-done
}
func say(s string, done chan<- struct{}) {
	defer close(done)
	for i := 0; i < 3; i++ {
		fmt.Println(s, i)
	}
}
func main() {
	if len(os.Args) == 3 {
		cell(os.Args[1], os.Args[2])
		return
	}
	fmt.Println("=== Ping: sleep in sender ===")
	ping(false)
	fmt.Println("=== Ping: sleep in receiver ===")
	ping(true)
	fmt.Println("=== Channel instead of WaitGroup ===")
	done := make(chan struct{})
	go say("goroutine", done)
	<-done
	fmt.Println("=== 12 table cells (child processes) ===")
	for _, state := range []string{"nil", "open", "closed"} {
		for _, op := range []string{"send", "receive", "close", "size"} {
			// Go's own deadlock detector reports the nil-block and open-receive cases.
			cmd := exec.Command(os.Args[0], state, op)
			out, err := cmd.CombinedOutput()
			fmt.Printf("%-6s %-7s => %s", state, op, out)
			if len(out) == 0 {
				fmt.Println(err)
			}
		}
	}
	fmt.Println("A blocked goroutine is not diagnosed as a global deadlock while other goroutines are active; it can leak indefinitely.")
}
