package main

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

var retainedSlices [][]byte
var retainedStrings []string

func readFile(i int) []byte {
	data := make([]byte, 1<<20)
	copy(data, fmt.Sprintf("HEADER-%04d", i))
	return data
}
func keepHeadersSlice(n int) [][]byte {
	h := make([][]byte, 0, n)
	for i := 0; i < n; i++ {
		d := readFile(i)
		h = append(h, d[:11])
	}
	return h
}
func keepHeadersCopy(n int) [][]byte {
	h := make([][]byte, 0, n)
	for i := 0; i < n; i++ {
		d := readFile(i)
		h = append(h, bytes.Clone(d[:11]))
	}
	return h
}
func keepStringsSlice(n int) []string {
	out := make([]string, 0, n)
	for i := 0; i < n; i++ {
		b := make([]byte, 1<<20)
		copy(b, fmt.Sprintf("HEAD-%04d", i))
		big := string(b)
		out = append(out, big[:5])
	}
	return out
}
func keepStringsClone(n int) []string {
	out := make([]string, 0, n)
	for i := 0; i < n; i++ {
		b := make([]byte, 1<<20)
		copy(b, fmt.Sprintf("HEAD-%04d", i))
		big := string(b)
		out = append(out, strings.Clone(big[:5]))
	}
	return out
}
func heap() uint64 {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.HeapAlloc
}
func main() {
	const n = 100
	baseline := heap()
	retainedSlices = keepHeadersSlice(n)
	a := heap()
	fmt.Printf("subslice retained: %.1f MiB\n", float64(a-baseline)/(1<<20))
	runtime.KeepAlive(retainedSlices)
	retainedSlices = nil
	baseline = heap()
	retainedSlices = keepHeadersCopy(n)
	a = heap()
	fmt.Printf("byte clones retained: %.3f MiB\n", float64(a-baseline)/(1<<20))
	runtime.KeepAlive(retainedSlices)
	retainedSlices = nil
	baseline = heap()
	retainedStrings = keepStringsSlice(n)
	a = heap()
	fmt.Printf("substring retained: %.1f MiB\n", float64(a-baseline)/(1<<20))
	runtime.KeepAlive(retainedStrings)
	retainedStrings = nil
	baseline = heap()
	retainedStrings = keepStringsClone(n)
	a = heap()
	fmt.Printf("string clones retained: %.3f MiB\n", float64(a-baseline)/(1<<20))
	runtime.KeepAlive(retainedStrings)
}
