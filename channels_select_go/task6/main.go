package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Result struct{ Word, Upper string }

var ErrQueueFull = errors.New("queue is full")

func producer(words []string, done <-chan struct{}) <-chan string {
	jobs := make(chan string)
	go func() {
		defer close(jobs)
		for _, w := range words {
			select {
			case jobs <- w:
			case <-done:
				return
			}
		}
	}()
	return jobs
}
func startWorkers(n int, jobs <-chan string, done <-chan struct{}) <-chan Result {
	results := make(chan Result)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				case w, ok := <-jobs:
					if !ok {
						return
					}
					r := Result{w, strings.ToUpper(w)}
					select {
					case results <- r:
					case <-done:
						return
					}
				}
			}
		}()
	}
	go func() { wg.Wait(); close(results) }()
	return results
}
func Submit(queue chan<- string, job string, timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case queue <- job:
		return nil
	case <-timer.C:
		return ErrQueueFull
	}
}

// Deliberately broken example: without done checks workers remain blocked on
// result sends when the consumer stops. Run in isolation: go run ./task6 leak
func demonstrateLeak() {
	jobs := make(chan string)
	results := make(chan Result)
	var wg sync.WaitGroup
	for w := 0; w < 3; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for word := range jobs {
				results <- Result{word, strings.ToUpper(word)}
			}
		}()
	}
	go func() {
		defer close(jobs)
		for i := 0; i < 100; i++ {
			jobs <- fmt.Sprintf("job-%d", i)
		}
	}()
	fmt.Println("got first result:", <-results)
	// The rest of the results are never consumed and the broken workers can't stop.
	time.Sleep(30 * time.Millisecond)
	fmt.Println("goroutines with missing done checks:", runtime.NumGoroutine(), "(expected >1: leaks)")
	// Do not call wg.Wait(): it would block indefinitely.
	_ = wg
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "leak" {
		demonstrateLeak()
		return
	}
	words := []string{"go", "channels", "select", "worker", "pool", "bounded", "cancellation", "hello", "world", "race"}
	done := make(chan struct{})
	for r := range startWorkers(3, producer(words, done), done) {
		fmt.Println(r)
	}
	close(done)
	done2 := make(chan struct{})
	results := startWorkers(3, producer(words, done2), done2)
	for i := 0; i < 3; i++ {
		r, ok := <-results
		if !ok {
			break
		}
		fmt.Println("early result:", r)
	}
	close(done2)
	for range results {
	} // wait for the closer to shut down
	time.Sleep(20 * time.Millisecond)
	fmt.Println("goroutines after cancellation:", runtime.NumGoroutine(), "(expected 1)")
	q := make(chan string, 1)
	fmt.Println("Submit success:", Submit(q, "job1", 10*time.Millisecond))
	fmt.Println("Submit full:", Submit(q, "job2", 10*time.Millisecond))
	// Demonstration of a leak, isolated in a child process: go run ./task6 leak
	// It keeps the blocked worker alive while main runs, rather than crashing the default run.
}
