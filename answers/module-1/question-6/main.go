package main

import "fmt"

// TODO: Implement this function.
func MarkAsPaid(isPaid *bool) {
	// Change the value that the pointer is pointing to.
	*isPaid = true
}

func main() {
	isUnpaid := false
	fmt.Printf("Status before: %v\n", isUnpaid) // Expected: Status before: false

	// Pass the memory address of isUnpaid to the function.
	MarkAsPaid(&isUnpaid)

	fmt.Printf("Status after: %v\n", isUnpaid) // Expected: Status after: true
}
