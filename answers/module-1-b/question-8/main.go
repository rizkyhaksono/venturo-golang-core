package main

import "fmt"

type ProductCategory uint8

const (
	_ ProductCategory = iota
	Goods
	Service
	Subscription
)

// TODO: Implement this function.
func GetCategoryType(category ProductCategory) string {
	// Use an if/else or switch statement.
	// If the category is Goods, return "Physical Item".
	// Otherwise, return "Digital Item or Service".
	switch category {
	case Goods:
		return "Physical Item"
	case Service:
		return "Digital Item or Service"
	default:
		return "Unknown Category"
	}
}

func main() {
	fmt.Printf("Goods: %s\n", GetCategoryType(Goods))     // Expected: Goods: Physical Item
	fmt.Printf("Service: %s\n", GetCategoryType(Service)) // Expected: Service: Digital Item or Service
}
