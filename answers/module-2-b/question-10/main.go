package main

import (
	"fmt"
	"sync"
	"time"
)

type Sale struct {
	Product string
	Amount  int
}

type Report struct {
	TotalSales int
	mu         sync.Mutex // To protect TotalSales
}

// updateReport listens on the channel and updates the report.
func (r *Report) updateReport(salesChannel chan Sale, wg *sync.WaitGroup) {
	defer wg.Done()
	for sale := range salesChannel {
		r.mu.Lock()
		fmt.Printf("Processing sale of %s for $%d\n", sale.Product, sale.Amount)
		r.TotalSales += sale.Amount
		r.mu.Unlock()
		time.Sleep(50 * time.Millisecond) // Simulate work
	}
}

func main() {
	var wg sync.WaitGroup
	report := Report{}
	salesChannel := make(chan Sale, 5) // A buffered channel

	wg.Add(1)
	go report.updateReport(salesChannel, &wg)

	// Simulate incoming sales from multiple sources.
	salesChannel <- Sale{Product: "Laptop", Amount: 1200}
	salesChannel <- Sale{Product: "Mouse", Amount: 25}
	salesChannel <- Sale{Product: "Keyboard", Amount: 75}
	salesChannel <- Sale{Product: "Monitor", Amount: 300}

	close(salesChannel) // Close the channel to signal no more sales are coming.
	wg.Wait()           // Wait for the report processor to finish.

	fmt.Printf("Final Report - Total Sales: $%d\n", report.TotalSales) // Expected: Final Report - Total Sales: $1600
}
