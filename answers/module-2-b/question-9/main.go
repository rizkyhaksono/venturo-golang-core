package main

import (
	"fmt"
	"time"
)

// This function starts a goroutine and immediately returns the channel.
func slowCalculation() chan int {
	resultChannel := make(chan int)
	go func() {
		time.Sleep(2 * time.Second)
		resultChannel <- 100 * 5 // The slow work
	}()
	return resultChannel
}

func main() {
	fmt.Println("Starting slow calculation...")
	resultCh := slowCalculation()

	// Wait for the result from the channel.
	result := <-resultCh

	fmt.Printf("Calculation finished. Result: %d\n", result)
}
