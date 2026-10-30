package main

import (
	"fmt"
	"strconv"
	"strings"
)

func JoinIDsSlow(ids []int) string {
	s := ""
	for i, id := range ids {
		if i > 0 {
			s += ","
		}
		s += fmt.Sprint(id)
	}
	return s
}
func JoinIDsFast(ids []int) string {
	var sb strings.Builder
	sb.Grow(len(ids) * 8)
	var tmp [20]byte
	for i, id := range ids {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.Write(strconv.AppendInt(tmp[:0], int64(id), 10))
	}
	return sb.String()
}
func main() { fmt.Println(JoinIDsSlow([]int{1, 2, 3})); fmt.Println(JoinIDsFast([]int{1, 2, 3})) }
