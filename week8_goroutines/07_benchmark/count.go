package main

import (
	"sync"
)

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for d := 2; d*d <= n; d++ {
		if n%d == 0 {
			return false
		}
	}
	return true
}

func CountSeq(n int) int {
	count := 0
	for value := 2; value < n; value++ {
		if isPrime(value) {
			count++
		}
	}
	return count
}

func CountPar(n, parts int) int {
	if parts < 1 {
		parts = 1
	}
	if parts > n-1 {
		parts = n - 1
	}
	results := make([]int, parts)
	var wg sync.WaitGroup
	chunk := (n - 1 + parts - 1) / parts
	for part := 0; part < parts; part++ {
		from := 2 + part*chunk
		to := from + chunk
		if to > n {
			to = n
		}
		if from >= n {
			break
		}
		wg.Add(1)
		go func(index, start, end int) {
			defer wg.Done()
			local := 0
			for value := start; value < end; value++ {
				if isPrime(value) {
					local++
				}
			}
			results[index] = local
		}(part, from, to)
	}
	wg.Wait()
	total := 0
	for _, value := range results {
		total += value
	}
	return total
}
