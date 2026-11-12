package main

import (
	"fmt"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
)

func busy(duration time.Duration) {
	end := time.Now().Add(duration)
	for time.Now().Before(end) { /* consume CPU */
	}
}

func runBusyWorkers() {
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); busy(2 * time.Second) }()
	}
	wg.Wait()
}

func preemptionDemo() {
	runtime.GOMAXPROCS(1)
	go func() {
		for { /* intentionally busy loop; runtime should preempt it */
		}
	}()
	time.Sleep(10 * time.Millisecond)
	fmt.Println("main is still alive")
}

func countThreads(label string, work func()) {
	before := pprof.Lookup("threadcreate").Count()
	work()
	after := pprof.Lookup("threadcreate").Count()
	fmt.Printf("%s: threadcreate profile count changed from %d to %d\n", label, before, after)
}

func main() {
	fmt.Println("CPU cores:", runtime.NumCPU(), "GOMAXPROCS:", runtime.GOMAXPROCS(0))
	fmt.Println("Running 8 CPU workers; run with GODEBUG=schedtrace=500 to see scheduler snapshots")
	runBusyWorkers()
	fmt.Println("Preemption demonstration:")
	preemptionDemo()
	fmt.Println("Thread profile demonstration (counts are cumulative):")
	countThreads("time.Sleep", func() {
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() { defer wg.Done(); time.Sleep(20 * time.Millisecond) }()
		}
		wg.Wait()
	})
}
