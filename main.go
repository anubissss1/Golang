package main

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"
)

func concatWithBuilder(values []string) string {
	var sb strings.Builder

	totalLength := 0
	for _, value := range values {
		totalLength += len(value)
	}

	sb.Grow(totalLength)

	for _, value := range values {
		sb.WriteString(value)
	}

	return sb.String()
}

func processPayloadOriginal(data []byte) string {
	trimmed := string(bytes.TrimSpace(data))
	clean := strings.ReplaceAll(trimmed, "\r", "")

	return clean
}

func processPayloadOptimized(data []byte) string {
	clean := bytes.TrimSpace(data)
	clean = bytes.ReplaceAll(clean, []byte("\r"), nil)

	return string(clean)
}

func main() {

	fmt.Println("===== STRINGS, BYTES AND RUNES =====")

	s := "Hello, 世界"

	fmt.Println("String:", s)
	fmt.Println("len(s) in bytes:", len(s))
	fmt.Println("Number of runes:", utf8.RuneCountInString(s))

	fmt.Println("First byte:", s[0])

	fmt.Println("\nIterating using range:")

	for i, r := range s {
		fmt.Printf("byte position %d: %c\n", i, r)
	}

	runes := []rune(s)

	fmt.Println("\nSecond rune:", string(runes[1]))

	// ==========================================
	// 2. TRIM FUNCTIONS
	// ==========================================

	fmt.Println("\n===== TRIM FUNCTIONS =====")

	text := "123oxo"

	fmt.Println("Original:", text)

	fmt.Println("TrimRight:",
		strings.TrimRight(text, "xo"))

	fmt.Println("TrimSuffix:",
		strings.TrimSuffix(text, "xo"))

	fmt.Println("TrimLeft:",
		strings.TrimLeft("xxxHello", "x"))

	fmt.Println("TrimPrefix:",
		strings.TrimPrefix("prefix-Hello", "prefix-"))

	// ==========================================
	// 3. STRING CONCATENATION
	// ==========================================

	fmt.Println("\n===== STRING CONCATENATION =====")

	values := []string{
		"Hello",
		" ",
		"Go",
		" ",
		"Developer",
	}

	result := concatWithBuilder(values)

	fmt.Println(result)

	// ==========================================
	// 4. []BYTE AND STRING
	// ==========================================

	fmt.Println("\n===== []BYTE AND STRING =====")

	data := []byte("  Hello World\r\n ")

	fmt.Println("Original bytes:", data)
	fmt.Println("Original string:", string(data))

	fmt.Println(
		"Original function:",
		processPayloadOriginal(data),
	)

	fmt.Println(
		"Optimized function:",
		processPayloadOptimized(data),
	)

	// ==========================================
	// 5. MAPS
	// ==========================================

	fmt.Println("\n===== MAPS =====")

	scores := map[string]int{
		"Alice": 90,
		"Bob":   0,
	}

	bobScore, bobExists := scores["Bob"]

	if bobExists {
		fmt.Println("Bob is in the map")
		fmt.Println("Bob's score:", bobScore)
	} else {
		fmt.Println("Bob is missing")
	}

	charlieScore, charlieExists := scores["Charlie"]

	if charlieExists {
		fmt.Println("Charlie is in the map")
		fmt.Println("Charlie's score:", charlieScore)
	} else {
		fmt.Println("Charlie is missing")
	}

	fmt.Println("Charlie's value:", charlieScore)

	// ==========================================
	// 6. INVENTORY
	// ==========================================

	fmt.Println("\n===== INVENTORY =====")

	inventory := map[string]int{
		"apples":  10,
		"bananas": 5,
	}

	// Read.
	fmt.Println("Apples:", inventory["apples"])

	// Update.
	inventory["bananas"] = 12

	// Insert.
	inventory["oranges"] = 8

	// Delete.
	delete(inventory, "apples")

	fmt.Println("Updated inventory:", inventory)
}
