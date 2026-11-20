# Task 7 — benchmarks

Run in the project root:

```sh
go test ./task7 -v
go test ./task7 -run '^$' -bench . -benchmem -cpu 1,2
```

`N = 10,000`, so the required sum is **49,995,000**.

To compute approximate cost per value, divide the `ns/op` of a sub-benchmark by 10,000. `B/op` is allocated bytes per full run, not size of the channel buffer alone. Performance depends on CPU, Go version, scheduler and memory allocator. Batches amortize synchronization but do not guarantee a fixed 10x improvement on all machines.
