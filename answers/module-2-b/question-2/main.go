package main

import (
	"fmt"
	"sync"
	"time"
)

func sayHelloWithWaitGroup(wg *sync.WaitGroup) {
	defer wg.Done() // This will be called when the function returns.
	time.Sleep(1 * time.Second)
	fmt.Println("Hello from goroutine!")
}

func main() {
	var wg sync.WaitGroup // Create a WaitGroup.

	fmt.Println("Hello from main!")

	wg.Add(1) // Increment the WaitGroup counter.
	go sayHelloWithWaitGroup(&wg)

	wg.Wait() // Block until the WaitGroup counter is zero.
	fmt.Println("Main function finished.")
}
