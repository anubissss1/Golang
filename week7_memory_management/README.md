# Week 7 — Memory Management (Go)

Complete runnable solutions for **WEEK 7. Memory Management.pdf**. Each folder is an independent Go `main` package. Requires Go 1.23+ for the project; the PDF was written with Go 1.24 and discusses Go 1.25–1.27, so escape-analysis output and benchmark numbers may differ.

## Quick start

Open a terminal **in this project folder**:

```bash
go version
go test ./...
go run ./02_escape
go run ./05_leaks
go run ./07_pool
go test -run '^$' -bench=. -benchmem ./03_allocations ./04_benchmarks ./07_pool
```

## Task 1 — Explore the stack (`01_stack`)

```bash
go run ./01_stack                     # depth(1_000_000)
go run ./01_stack 10000              # quicker test
# In 01_stack/main.go, change [64]byte to [1024]byte and rerun.
go build -gcflags='-S -l' -o /tmp/stackdemo ./01_stack 2> stack_assembly.txt
# Search stack_assembly.txt for 'main.depth STEXT' and 'TEXT main.depth(SB)'.
```

With the `[64]byte` array, a stack frame needs at least the 64-byte array plus saved return/call data and other temporaries; the **actual** frame size is architecture/compiler-dependent. With `[1024]byte`, each frame is roughly 960 bytes larger, so one million calls can approach or exceed the goroutine's 1 GB stack limit. Do not assume the larger version succeeds. The code uses `pad` after the recursive call so the compiler cannot simply discard it. Remove `+ int(pad[(n+1)%len(pad)])` and compile again: if `pad` is otherwise not observable, dead-code elimination can remove the array and dramatically shrink the frame.

**Stack overflow demonstration:** `go run ./01_stack forever` deliberately terminates with a **fatal runtime stack overflow** (not a recoverable panic). Run this separately; it is expected to fail and cannot be caught with `recover`. The stack grows by allocating a larger contiguous stack and copying the old frames. The Go runtime adjusts pointers into the moved stack, making growth safe.

## Task 2 — Predict and check (`02_escape`)

```bash
go run ./02_escape
go build -gcflags='-m -l' -o /tmp/escape_demo ./02_escape 2> escape_report.txt
```

| Function | Prediction / reason | Observed with Go 1.23.2 (`-m -l`) |
|---|---|---|
| `f1` | Local `x` stays on stack; return integer by value | No `x` heap report |
| `f2` | `x` escapes; its address is returned | `moved to heap: x` |
| `f3` | `p` itself does not escape; only reads fields | `p does not escape` |
| `f4` | `p` passed to `fmt.Println` is boxed; may escape | `p escapes to heap` |
| `f5` | Backing array escapes; returned slice refers to it | `make([]int, 10) escapes to heap` |
| `f6` | Backing array need not escape; only `len` used | `make([]int, 10) does not escape` |
| `f7` | Runtime-sized backing array may escape | `make([]int, n) escapes to heap` |
| `f8` | Captured `x` escapes with returned closure | `moved to heap: x`; function literal escapes |
| `f9` | Pointer `p` leaks into global; caller's object escapes | `leaking param: p`; `main`'s `p` moved to heap |
| `f10` | Large local array: may be heap allocated or optimized away | No report for `big` in this build: compiler can eliminate the array because only zero-valued `big[0]` is read |

**Important:** `f4` prints a struct and may allocate due to interface boxing; `f10` is a case where the source-level array doesn't need storage at all. This is why compiler reports must be checked, not assumed. Compiler decisions change with optimization and Go versions.

## Task 3 — Remove allocations (`03_allocations`)

`JoinIDsSlow` repeatedly concatenates strings and uses `fmt.Sprint(id)`. `JoinIDsFast` uses a `strings.Builder` with `Grow` and `strconv.AppendInt` into a stack-local temporary buffer. Both return the same comma-separated text.

```bash
go test -v ./03_allocations
go test -run '^$' -bench=. -benchmem ./03_allocations
go build -gcflags='-m -l' -o /tmp/join_demo ./03_allocations 2> join_escape.txt
```

The `fmt.Sprint(id)` line can make `id` escape via interface boxing. The fast version typically has **one** allocation: the returned string's backing storage must survive the function. `AllocsPerRun` checks that the fast version uses at most one allocation. The exact slow count depends on the values and compiler; it is not necessarily the PDF's example count.

## Task 4 — Benchmarks (`04_benchmarks`)

```bash
go test -run '^$' -bench=. -benchmem ./04_benchmarks
go build -gcflags='-m -l' -o /tmp/point_demo ./04_benchmarks 2> point_escape.txt
```

`//go:noinline` prevents the compiler from inlining constructors and obscuring the comparison. `NewPointValue` returns a small value and usually needs **0 heap allocations**. `NewPointPointer` returns a pointer to a new value, usually **1 allocation of 16 bytes**. `CollectNoCap(1000)` repeatedly grows and copies backing arrays, usually about **12 allocations / 25 KB**, while `CollectWithCap(1000)` allocates one **8 KB** array. The extra bytes are abandoned intermediate backing arrays from geometric capacity growth.

## Task 5 — Find the leak (`05_leaks`)

```bash
go run ./05_leaks
```

`keepHeadersSlice` keeps `data[:11]`, a slice header pointing into the original 1 MiB backing array. `keepHeadersCopy` uses `bytes.Clone(data[:11])` to retain only 11 bytes. Similarly, `big[:5]` keeps a large string's backing bytes reachable; `strings.Clone(big[:5])` copies just five bytes. The demo uses `runtime.GC`, `runtime.ReadMemStats`, and `runtime.KeepAlive` so the references are still live at measurement time. Results are approximate because the runtime also allocates memory.

## Task 6 — Watch the GC (`06_gc`)

The program allocates 80,000 live `Record`s and 10 million short-lived records. `consume` writes the last record to a global so the temporary objects really **escape**; otherwise the compiler could optimize away the intended heap allocations.

```bash
go build -o gc_demo ./06_gc
GOGC=50 ./gc_demo
GOGC=100 ./gc_demo
GOGC=200 ./gc_demo
GOGC=400 ./gc_demo
GOGC=off ./gc_demo
GOGC=off GOMEMLIMIT=100MiB ./gc_demo
GODEBUG=gctrace=1 GOGC=100 ./gc_demo
```

On **Windows PowerShell**, set environment variables for one run with `$env:GOGC='50'; .\gc_demo.exe`; on **macOS/Linux**, use the commands above. Avoid `GOGC=off` without a memory limit on memory-constrained machines: this example can use several gigabytes.

Sample `gctrace` format: `gc 61 @0.453s 6%: 0.012+1.9+0.048 ms clock, ..., 46->51->27 MB, 54 MB goal`. Here `61` is the cycle number, `@0.453s` is elapsed time, `6%` is approximate GC CPU percentage, three timing numbers are pause/mark/pause phases, `46->51->27` shows heap sizes before/during/after GC, and `54 MB goal` is the next heap goal.

**Choosing settings:** `GOGC=50` trades more frequent GC for lower memory; `GOGC=200` can suit a batch job with more RAM. `GOGC=off` disables GC triggered by heap growth, but a `GOMEMLIMIT` can still trigger collections. `GOMEMLIMIT` is a **soft** limit, not a hard cap.

## Task 7 — Reuse buffers (`07_pool`)

```bash
go test -v ./07_pool
go test -run '^$' -bench=. -benchmem ./07_pool
```

Both `RenderNew` and `RenderPool` write `id=<n>\n` for 200 IDs and return the byte count. `RenderNew` allocates a new buffer's backing storage; `RenderPool` borrows a buffer from `sync.Pool` and may reuse its backing storage. `Reset` clears old content before reuse. **Never use a buffer after `Put`**: another goroutine may borrow it, causing data races or corruption. A pool is a cache; GC may discard its contents, so occasional allocations are expected.

## Bonus — Stack allocation threshold (`08_bonus`)

```bash
go build -gcflags='-m -l' -o /tmp/bonus_demo ./08_bonus 2>&1
```

On Go 1.23.2, `make([]int, 1000)` **does not escape**, while `make([]int, 10000)` **escapes to heap**. For this version's constant-size `make` rule, the threshold is about **64 KiB** of backing storage: on 64-bit systems, 8192 `int`s occupy 65,536 bytes, and 8193 occupy 65,544 bytes. Test the boundary in your Go version; implementation thresholds may change. The compiler limits very large stack objects to avoid costly stack growth and copying.

## Measured benchmark results (example environment only)

Measurements below were taken with **Go 1.23.2, Linux/amd64**, using `-benchtime=100ms`. Re-run on your own laptop for the defence.

| Benchmark | ns/op | B/op | allocs/op |
|---|---:|---:|---:|
| JoinIDsSlow | 2566 | 240 | 19 |
| JoinIDsFast | 325.5 | 80 | 1 |
| NewPointValue | 2.368 | 0 | 0 |
| NewPointPointer | 43.51 | 16 | 1 |
| CollectNoCap | 18492 | 25208 | 12 |
| CollectWithCap | 7346 | 8192 | 1 |
| RenderNew | 8489 | 4096 | 1 |
| RenderPool | 5664 | 0 | 0 |

`05_leaks` on the same environment retained ~100 MiB with subslices and substrings, versus ~0.009 MiB and ~0.004 MiB for clones. The GC demonstration measured **118 GC cycles** with `GOGC=100` and **46 cycles** with `GOGC=off GOMEMLIMIT=100MiB` on one run. These are *not* universal expected results.

## Defence questions — short answers

1. **Stack vs heap?** Each goroutine has its own stack for call frames and short-lived values; heap memory is shared and GC-managed. Stack allocation is generally cheaper; heap allocation adds GC work.
2. **Who decides? Does `new` or `&` force the heap?** The Go compiler's escape analysis decides. Neither `new` nor `&` automatically means heap.
3. **What is escape analysis? Four reasons?** It determines whether values can safely remain on the stack. Values may escape when returned by pointer, stored globally, captured by a long-lived closure, or stored in an escaping interface. Large/dynamic allocations are other reasons.
4. **How to inspect? `leaking param`?** `go build -gcflags='-m -l' .` prints decisions. `leaking param: p` means the pointer argument can flow somewhere beyond the function, potentially forcing its referent onto the heap.
5. **Count allocations?** `testing.AllocsPerRun` in tests; `go test -bench=. -benchmem` for `allocs/op` and `B/op` in benchmarks.
6. **How does GC free memory? Roots?** The concurrent mark-and-sweep collector traces pointers starting from roots (goroutine stacks, globals, runtime roots); unreachable objects can be reclaimed.
7. **`GOGC` and `GOMEMLIMIT`?** `GOGC` adjusts heap-growth targets (memory/CPU trade-off); `GOMEMLIMIT` sets a soft overall runtime memory limit that can increase GC frequency.
8. **Why can a tiny slice/substring keep large memory?** It still references the original large backing array/string. `bytes.Clone` and `strings.Clone` create independent small copies.

### Extra questions from the PDF

- **Why is `return &u` safe?** Escape analysis puts `u` in memory that remains valid after return, typically the heap.
- **Why does the stack copy safely?** The runtime updates pointers that refer into the old stack when it moves stack frames.
- **Why is stack overflow not recoverable?** It is a fatal runtime error, not a regular `panic`.
- **Why is `JoinIDsFast` still 1 allocation?** The returned string needs backing bytes that survive the function.
- **Why exactly 12 `CollectNoCap` allocations here?** The Go runtime repeatedly grows slice capacity; for this Go version and `int` element size, it took 12 backing arrays to reach 1000 elements. The growth algorithm is an implementation detail.
- **Why can `GOGC=off` be slower?** Large heaps put pressure on RAM, memory bandwidth, caches and the OS; turning off GC is not automatically faster.
- **Why must `sync.Pool` not be permanent storage?** GC may clear pooled values at any time.
