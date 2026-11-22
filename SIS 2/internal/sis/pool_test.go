package sis

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestAnalyzeTable(t *testing.T) {
	tests := []struct {
		input                string
		words, chars, unique int
		longest              string
		invalid              bool
	}{
		{"Hello Go world", 3, 14, 3, "Hello", false},
		{"Go, go! 42", 3, 10, 2, "Go", false},
		{"  \t  ", 0, 0, 0, "", true},
		{"é à", 2, 3, 2, "é", false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%q", tt.input), func(t *testing.T) {
			got, err := Analyze(tt.input)
			if tt.invalid {
				if !errors.Is(err, ErrInvalidText) {
					t.Fatalf("want invalid, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got.Words != tt.words || got.Characters != tt.chars || got.UniqueWords != tt.unique || got.LongestWord != tt.longest {
				t.Fatalf("unexpected result %+v", got)
			}
		})
	}
}
func TestCreateRetrieveValidateAndSyncProcessing(t *testing.T) {
	p, _ := NewPool(2, 2)
	defer p.Shutdown()
	if _, err := p.Create(" "); !errors.Is(err, ErrInvalidText) {
		t.Fatalf("invalid input: %v", err)
	}
	if _, err := p.Get(9999); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing task: %v", err)
	}
	a, err := p.Create("hello go")
	if err != nil {
		t.Fatal(err)
	}
	b, err := p.Create("three little words")
	if err != nil {
		t.Fatal(err)
	}
	if a.ID == b.ID || a.State != Pending {
		t.Fatal("invalid identities or states")
	}
	if err := p.ProcessSync(a.ID); err != nil {
		t.Fatal(err)
	}
	if err := p.ProcessSync(a.ID); !errors.Is(err, ErrInvalidState) {
		t.Fatalf("duplicate sync processing: %v", err)
	}
	got, _ := p.Get(a.ID)
	if got.State != Completed || got.Result.Words != 2 || got.Result.Characters != 8 {
		t.Fatalf("sync result %+v", got)
	}
	other, _ := p.Get(b.ID)
	if other.State != Pending || other.Result.Words != 0 {
		t.Fatalf("different task overwritten: %+v", other)
	}
	if err := p.Submit(a.ID); !errors.Is(err, ErrInvalidState) {
		t.Fatalf("resubmit completed: %v", err)
	}
}
func TestWorkerPoolProcessingAndShutdown(t *testing.T) {
	p, _ := NewPool(3, 5)
	if p.WorkerCount() != 3 || p.QueueCapacity() != 5 {
		t.Fatal("wrong configuration")
	}
	ids := make([]uint64, 0, 80)
	for i := 0; i < 80; i++ {
		task, err := p.Create(fmt.Sprintf("hello task%d world", i))
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, task.ID)
		for {
			err = p.Submit(task.ID)
			if err == nil {
				break
			}
			if !errors.Is(err, ErrQueueFull) {
				t.Fatal(err)
			} // yield until a worker takes a task
			select {
			case <-time.After(time.Millisecond):
			case <-time.After(10 * time.Millisecond):
				t.Fatal("queue did not drain")
			}
		}
	}
	p.Shutdown()
	p.Shutdown()
	for _, id := range ids {
		task, err := p.Get(id)
		if err != nil || task.State != Completed || task.Result.Words != 3 {
			t.Fatalf("task %d: %+v, %v", id, task, err)
		}
	}
	if _, err := p.Create("stopped"); !errors.Is(err, ErrStopped) {
		t.Fatalf("create after stop: %v", err)
	}
	if err := p.Submit(ids[0]); !errors.Is(err, ErrStopped) {
		t.Fatalf("submit after stop: %v", err)
	}
}
func TestActualConcurrentOverlap(t *testing.T) {
	entered := make(chan struct{}, 2)
	release := make(chan struct{})
	fn := func(s string) (Analysis, error) { entered <- struct{}{}; <-release; return Analyze(s) }
	p, _ := NewPoolWithAnalyzer(2, 2, fn)
	defer func() { close(release); p.Shutdown() }()
	for i := 0; i < 2; i++ {
		v, _ := p.Create("worker overlap")
		if err := p.Submit(v.ID); err != nil {
			t.Fatal(err)
		}
	}
	for i := 0; i < 2; i++ {
		select {
		case <-entered:
		case <-time.After(2 * time.Second):
			t.Fatal("2 workers did not overlap")
		}
	}
}
func TestQueueFullAndDuplicateSubmission(t *testing.T) {
	entered := make(chan struct{}, 2)
	release := make(chan struct{})
	fn := func(s string) (Analysis, error) { entered <- struct{}{}; <-release; return Analyze(s) }
	p, _ := NewPoolWithAnalyzer(2, 1, fn)
	defer func() { close(release); p.Shutdown() }()
	for i := 0; i < 2; i++ {
		v, _ := p.Create("busy worker")
		if err := p.Submit(v.ID); err != nil {
			t.Fatal(err)
		}
		select {
		case <-entered:
		case <-time.After(time.Second):
			t.Fatal("workers did not start")
		}
	}
	q, _ := p.Create("queued")
	if err := p.Submit(q.ID); err != nil {
		t.Fatal(err)
	}
	if err := p.Submit(q.ID); !errors.Is(err, ErrInvalidState) {
		t.Fatalf("duplicate: %v", err)
	}
	extra, _ := p.Create("retry later")
	if err := p.Submit(extra.ID); !errors.Is(err, ErrQueueFull) {
		t.Fatalf("full queue: %v", err)
	}
	snap, _ := p.Get(extra.ID)
	if snap.State != Pending {
		t.Fatalf("rejected task state=%s", snap.State)
	}
}
func TestConcurrentReadsWritesAndShutdownRace(t *testing.T) {
	p, _ := NewPool(4, 16)
	var wg sync.WaitGroup
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				v, err := p.Create(fmt.Sprintf("task %d %d", i, j))
				if errors.Is(err, ErrStopped) {
					return
				}
				if err != nil {
					t.Error(err)
					return
				}
				_ = p.GetOrZero(v.ID)
				p.List()
				_ = p.Submit(v.ID)
			}
		}(i)
	}
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				p.List()
				_ = p.QueueLength()
			}
		}()
	}
	wg.Wait()
	var stops sync.WaitGroup
	for i := 0; i < 5; i++ {
		stops.Add(1)
		go func() { defer stops.Done(); p.Shutdown() }()
	}
	stops.Wait()
	for _, v := range p.List() {
		if v.State == Queued || v.State == Processing {
			t.Fatalf("accepted task %d unfinished", v.ID)
		}
	}
}
func TestFailureAndConfiguration(t *testing.T) {
	if _, err := NewPool(1, 2); !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("worker count %v", err)
	}
	if _, err := NewPool(2, 0); !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("queue capacity %v", err)
	}
	if _, err := NewPoolWithAnalyzer(2, 2, nil); err == nil {
		t.Fatal("nil analyzer accepted")
	}
	p, _ := NewPoolWithAnalyzer(2, 2, func(string) (Analysis, error) { return Analysis{}, errors.New("analysis failed") })
	v, _ := p.Create("valid input")
	if err := p.Submit(v.ID); err != nil {
		t.Fatal(err)
	}
	p.Shutdown()
	v, _ = p.Get(v.ID)
	if v.State != Failed || v.Error != "analysis failed" {
		t.Fatalf("failure result %+v", v)
	}
}
func TestSynchronousAndAsynchronousDoNotDoubleProcess(t *testing.T) {
	var calls atomic.Int64
	p, _ := NewPoolWithAnalyzer(2, 2, func(s string) (Analysis, error) { calls.Add(1); return Analyze(s) })
	defer p.Shutdown()
	v, _ := p.Create("one job")
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); _ = p.ProcessSync(v.ID) }()
	go func() { defer wg.Done(); _ = p.Submit(v.ID) }()
	wg.Wait()
	p.Shutdown()
	if calls.Load() != 1 {
		t.Fatalf("task processed %d times", calls.Load())
	}
}
