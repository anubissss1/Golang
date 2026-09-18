package main

import "fmt"

func main() {
	/// #1
	numbers := [5]float64{1.1, 2.2, 3.3, 4.4, 5.5}
	fmt.Println(numbers)

	/// #2
	var values [5]int = [5]int{10, 20, 30}
	fmt.Println(values)

	/// #3
	var matrix [2][3]int = [2][3]int{{1, 2, 3}, {4, 5, 6}}
	for _, row := range matrix {
		for _, value := range row {
			fmt.Print(value, " ")
		}
		fmt.Println()
	}
}
