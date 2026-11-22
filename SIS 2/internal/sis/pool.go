// Package sis implements a bounded, concurrency-safe text-analysis worker pool.
package sis

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

var (
	ErrInvalidText   = errors.New("text must contain at least one non-whitespace character and be at most 10000 bytes")
	ErrNotFound      = errors.New("task not found")
	ErrInvalidState  = errors.New("task must be pending to be submitted")
	ErrQueueFull     = errors.New("task queue is full")
	ErrStopped       = errors.New("worker pool is stopped")
	ErrInvalidConfig = errors.New("workers must be at least 2 and queue capacity at least 1")
)

type State string

const (
	Pending    State = "Pending"
	Queued     State = "Queued"
	Processing State = "Processing"
	Completed  State = "Completed"
	Failed     State = "Failed"
)

// Analysis is a domain-specific computation result.
type Analysis struct {
	Uppercase   string
	Words       int
	Characters  int
	UniqueWords int
	LongestWord string
}

// Task is a value snapshot: callers never receive pointers into shared state.
type Task struct {
	ID     uint64
	Text   string
	State  State
	Result Analysis
	Error  string
}

// Processor abstracts domain-specific logic (Checkpoint 1 interface requirement).
// Tests can inject a processor to control worker timing without arbitrary sleeps.
type Processor interface {
	Process(string) (Analysis, error)
}

// TextProcessor is the production implementation of Processor.
type TextProcessor struct{}

func (TextProcessor) Process(s string) (Analysis, error) { return Analyze(s) }

// Analyzer adapts a function to the Processor interface.
type Analyzer func(string) (Analysis, error)

func (a Analyzer) Process(s string) (Analysis, error) { return a(s) }

// Analyze computes useful statistics without mutating the input.
func Analyze(input string) (Analysis, error) {
	if strings.TrimSpace(input) == "" || len(input) > 10000 {
		return Analysis{}, ErrInvalidText
	}
	words := strings.FieldsFunc(input, func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsNumber(r) })
	seen := make(map[string]struct{})
	longest := ""
	for _, w := range words {
		seen[strings.ToLower(w)] = struct{}{}
		if utf8.RuneCountInString(w) > utf8.RuneCountInString(longest) {
			longest = w
		}
	}
	return Analysis{Uppercase: strings.ToUpper(input), Words: len(words), Characters: utf8.RuneCountInString(input), UniqueWords: len(seen), LongestWord: longest}, nil
}

// Pool owns its queue and task map. mu protects all map contents, nextID,
// stopped, and the send/close lifecycle of the queue.
type Pool struct {
	mu          sync.RWMutex
	tasks       map[uint64]Task
	nextID      uint64
	jobs        chan uint64
	stopped     bool
	wg          sync.WaitGroup
	analyzer    Processor
	workerCount int
}

func NewPool(workers, queueCapacity int) (*Pool, error) {
	return NewPoolWithProcessor(workers, queueCapacity, TextProcessor{})
}

func NewPoolWithAnalyzer(workers, queueCapacity int, fn Analyzer) (*Pool, error) {
	if fn == nil {
		return nil, errors.New("analyzer cannot be nil")
	}
	return NewPoolWithProcessor(workers, queueCapacity, fn)
}

func NewPoolWithProcessor(workers, queueCapacity int, fn Processor) (*Pool, error) {
	if workers < 2 || queueCapacity < 1 {
		return nil, ErrInvalidConfig
	}
	if fn == nil {
		return nil, errors.New("analyzer cannot be nil")
	}
	p := &Pool{tasks: make(map[uint64]Task), jobs: make(chan uint64, queueCapacity), analyzer: fn, workerCount: workers}
	p.wg.Add(workers)
	for i := 0; i < workers; i++ {
		go p.worker()
	}
	return p, nil
}

func (p *Pool) Create(text string) (Task, error) {
	if strings.TrimSpace(text) == "" || len(text) > 10000 {
		return Task{}, ErrInvalidText
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.stopped {
		return Task{}, ErrStopped
	}
	p.nextID++
	task := Task{ID: p.nextID, Text: text, State: Pending}
	p.tasks[task.ID] = task
	return task, nil
}

// Submit is non-blocking. Sending with mu held makes channel close safe:
// Shutdown cannot close the queue until this send has finished.
// ProcessSync preserves the synchronous processing operation from Checkpoint 1.
// The state claim is atomic with Submit so a task can never be processed twice.
func (p *Pool) ProcessSync(id uint64) error {
	p.mu.Lock()
	if p.stopped {
		p.mu.Unlock()
		return ErrStopped
	}
	t, ok := p.tasks[id]
	if !ok {
		p.mu.Unlock()
		return ErrNotFound
	}
	if t.State != Pending {
		p.mu.Unlock()
		return fmt.Errorf("%w: task %d is %s", ErrInvalidState, id, t.State)
	}
	t.State = Processing
	p.tasks[id] = t
	p.mu.Unlock()
	result, err := p.analyzer.Process(t.Text)
	p.mu.Lock()
	t = p.tasks[id]
	if err != nil {
		t.State = Failed
		t.Error = err.Error()
	} else {
		t.State = Completed
		t.Result = result
	}
	p.tasks[id] = t
	p.mu.Unlock()
	return err
}

func (p *Pool) Submit(id uint64) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.stopped {
		return ErrStopped
	}
	t, ok := p.tasks[id]
	if !ok {
		return ErrNotFound
	}
	if t.State != Pending {
		return fmt.Errorf("%w: task %d is %s", ErrInvalidState, id, t.State)
	}
	select {
	case p.jobs <- id:
		t.State = Queued
		p.tasks[id] = t
		return nil
	default:
		// Rejected tasks stay Pending, allowing a later submission.
		return ErrQueueFull
	}
}

// GetOrZero is a convenience for concurrent monitoring; missing IDs return zero value.
func (p *Pool) GetOrZero(id uint64) Task { t, _ := p.Get(id); return t }

func (p *Pool) Get(id uint64) (Task, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	t, ok := p.tasks[id]
	if !ok {
		return Task{}, ErrNotFound
	}
	return t, nil
}

func (p *Pool) List() []Task {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]Task, 0, len(p.tasks))
	for _, t := range p.tasks {
		out = append(out, t)
	}
	return out
}

func (p *Pool) WorkerCount() int   { return p.workerCount }
func (p *Pool) QueueCapacity() int { return cap(p.jobs) }
func (p *Pool) QueueLength() int   { return len(p.jobs) } // diagnostic snapshot only

func (p *Pool) worker() {
	defer p.wg.Done()
	for id := range p.jobs {
		p.mu.Lock()
		t := p.tasks[id]
		t.State = Processing
		p.tasks[id] = t
		p.mu.Unlock()

		result, err := p.analyzer.Process(t.Text) // never hold mu during long computation

		p.mu.Lock()
		t = p.tasks[id]
		if err != nil {
			t.State = Failed
			t.Error = err.Error()
		} else {
			t.State = Completed
			t.Result = result
		}
		p.tasks[id] = t
		p.mu.Unlock()
	}
}

// Shutdown stops new submissions, closes the queue once, drains accepted work,
// and waits for all workers. Multiple concurrent callers are safe.
func (p *Pool) Shutdown() {
	p.mu.Lock()
	if !p.stopped {
		p.stopped = true
		close(p.jobs)
	}
	p.mu.Unlock()
	p.wg.Wait()
}
