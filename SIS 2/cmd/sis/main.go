package main

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"sis-cp2/internal/sis"
)

func main() {
	pool, err := sis.NewPool(3, 5)
	if err != nil {
		log.Fatal(err)
	}
	samples := []string{
		"Go channels connect concurrent workers", "The quick brown fox jumps over the lazy dog",
		"Concurrency requires careful synchronization", "Go Go GO channels channels",
		"A worker pool limits the number of goroutines", "Race detectors help catch shared-memory bugs",
		"Bounded queues prevent unlimited growth", "Graceful shutdown processes all accepted tasks",
	}
	ids := make([]uint64, 0, len(samples))
	for _, s := range samples {
		t, err := pool.Create(s)
		if err != nil {
			log.Fatal(err)
		}
		ids = append(ids, t.ID)
		// This retry is only for the demonstration; Submit itself never blocks.
		for {
			err = pool.Submit(t.ID)
			if err == nil {
				break
			}
			if !errors.Is(err, sis.ErrQueueFull) {
				log.Fatal(err)
			}
			time.Sleep(time.Millisecond)
		}
	}
	// Shutdown drains the queue and waits for all results to be saved.
	pool.Shutdown()
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	fmt.Printf("Workers: %d, queue capacity: %d\n", pool.WorkerCount(), pool.QueueCapacity())
	for _, id := range ids {
		t, err := pool.Get(id)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Task %d | %-9s | words=%d unique=%d chars=%d longest=%q\n", t.ID, t.State, t.Result.Words, t.Result.UniqueWords, t.Result.Characters, t.Result.LongestWord)
	}
}
