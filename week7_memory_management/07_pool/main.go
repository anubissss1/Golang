package main

import (
	"bytes"
	"fmt"
	"strconv"
	"sync"
)

var bufPool = sync.Pool{New: func() any { b := new(bytes.Buffer); b.Grow(4096); return b }}

func writeIDs(buf *bytes.Buffer, ids []int) {
	var tmp [20]byte
	for _, id := range ids {
		buf.WriteString("id=")
		buf.Write(strconv.AppendInt(tmp[:0], int64(id), 10))
		buf.WriteByte('\n')
	}
}
func RenderNew(ids []int) int {
	buf := new(bytes.Buffer)
	buf.Grow(4096)
	writeIDs(buf, ids)
	return buf.Len()
}
func RenderPool(ids []int) int {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	writeIDs(buf, ids)
	n := buf.Len()
	bufPool.Put(buf)
	return n
}
func main() {
	ids := make([]int, 200)
	for i := range ids {
		ids[i] = i
	}
	fmt.Println("new:", RenderNew(ids), "pooled:", RenderPool(ids))
}
