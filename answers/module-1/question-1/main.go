package main

import "fmt"

func main() {
	// TODO: Declare variables with the correct types.
	var age int8
	var population int32
	var worldPopulation int64

	// Assign values
	age = 35
	population = 850000
	worldPopulation = 8100000000

	// The fmt.Printf function's %T verb prints the type of the variable.
	fmt.Printf("Age: %d (Type: %T)\n", age)                          // Expected: Age: 35 (Type: int8)
	fmt.Printf("Population: %d (Type: %T)\n", population)            // Expected: Population: 850000 (Type: int32)
	fmt.Printf("World Population: %d (Type: %T)\n", worldPopulation) // Expected: World Population: 8100000000 (Type: int64)
}
