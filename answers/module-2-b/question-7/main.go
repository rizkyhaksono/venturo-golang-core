package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	jsonString := `{"Laptops": 25, "Keyboards": 100, "Mice": 150}`
	var productStock map[string]int

	err := json.Unmarshal([]byte(jsonString), &productStock)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Stock for Keyboards: %d\n", productStock["Keyboards"]) // Expected: Stock for Keyboards: 100
}
