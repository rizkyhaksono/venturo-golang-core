package main

import "fmt"

// TODO: Define Notifier interface and ConsoleNotifier struct
type Notifier interface {
	Send(message string)
}

type ConsoleNotifier struct{}

// Implement the Send method for ConsoleNotifier
func (c ConsoleNotifier) Send(message string) {
	fmt.Printf("[CONSOLE NOTIFIER]: %s\n", message)
}

// TODO: Define OrderService struct
// It should hold a Notifier.
type OrderService struct {
	notifier Notifier
}

// NewOrderService is the constructor that "injects" the dependency.
func NewOrderService(notifier Notifier) *OrderService {
	return &OrderService{
		notifier: notifier,
	}
}

// PlaceOrder should use the injected notifier.
func (s *OrderService) PlaceOrder() {
	s.notifier.Send("Your order has been placed!")
}

func main() {
	// Create the concrete dependency
	consoleNotifier := ConsoleNotifier{}
	// Inject it into the service
	orderService := NewOrderService(consoleNotifier)

	orderService.PlaceOrder() // Expected: [CONSOLE NOTIFIER]: Your order has been placed!
}
