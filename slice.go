package main

import "fmt"

func main() {
	// #1
	var slice1 = []int{1, 2, 3, 4, 5}

	var slice2 = make([]int, 5)

	var slice3 []int

	var slice4 = make([]int, 0, 5)

	fmt.Println("4 different ways of slice:")
	fmt.Println(slice1, slice2, slice3, slice4)

	// #2
	cars := []string{"Ferrari", "Honda", "Mercedes", "Chevrolet"}

	fmt.Println("Cars slice:")
	fmt.Println(cars)
	fmt.Println("cars length:", len(cars))
	fmt.Println("cars capacity:", cap(cars))

	cars = append(cars, "Toyota", "BMW")
	fmt.Println("Cars slice after append:")
	fmt.Println(cars)
	fmt.Println("cars length:", len(cars))
	fmt.Println("cars capacity:", cap(cars))

	// #3
	original := [5]int{10, 20, 30, 40, 50}

	var sliceA = original[1:4]
	fmt.Println("Original array:", original)
	fmt.Println("Slice A:", sliceA)

	sliceA[0] = 100

	fmt.Println("Original array after modifying slice A:", original)
	fmt.Println("Slice A after modification:", sliceA)

	original[2] = 200
	fmt.Println("Original array after modifying original:", original)
	fmt.Println("Slice A after modifying original:", sliceA)
}
