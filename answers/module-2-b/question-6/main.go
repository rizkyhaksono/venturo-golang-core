package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	categoryCount := map[string]int{
		"Goods":   150,
		"Service": 80,
	}

	jsonData, err := json.Marshal(categoryCount)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonData)) // Expected: {"Goods":150,"Service":80}
}
