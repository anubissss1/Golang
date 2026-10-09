package main

import "fmt"

// IPAddr is an IPv4 address represented as four bytes.
type IPAddr [4]byte

func (ip IPAddr) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

type Book struct {
	Title  string
	Author string
}

func (b *Book) String() string {
	return b.Title + " by " + b.Author
}

func main() {
	hosts := map[string]IPAddr{
		"loopback":  {127, 0, 0, 1},
		"googleDNS": {8, 8, 8, 8},
	}
	for name, ip := range hosts {
		fmt.Printf("%v: %v\n", name, ip)
	}

	book := Book{Title: "The Go Programming Language", Author: "Donovan & Kernighan"}
	fmt.Println(book)  // Book value lacks String() in its method set: default struct format.
	fmt.Println(&book) // *Book has String(): custom string format.
}
