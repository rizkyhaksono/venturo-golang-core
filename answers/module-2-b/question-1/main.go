package main

import (
	"fmt"
	"time"
)

func sayHelloAsync() {
	time.Sleep(1 * time.Second)
	fmt.Println("Hello from goroutine!")
}

func main() {
	fmt.Println("Hello from main!")
	go sayHelloAsync()
	// What happens if you add a time.Sleep(2 * time.Second) here?
	time.Sleep(2 * time.Second)
	fmt.Println("Bye from main!")
	// If you add time.Sleep(2 * time.Second) here, the main function will wait for 2 seconds before exiting.
	// This allows the goroutine to complete its execution before the main function exits.
}
