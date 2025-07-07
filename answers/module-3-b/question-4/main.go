// ... (reuse code from Challenge 3) ...
package main

import "fmt"

type Notifier interface {
	Send(message string)
}

type ConsoleNotifier struct{}

func (c ConsoleNotifier) Send(message string) {
	fmt.Printf("[CONSOLE NOTIFIER]: %s\n", message)
}

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

// TODO: Define EmailNotifier struct and implement the Send method
type EmailNotifier struct{}

func (e EmailNotifier) Send(message string) {
	fmt.Printf("[CONSOLE NOTIFIER]: %s\n", message)
}

func main() {
	// --- Scenario 1 ---
	consoleNotifier := ConsoleNotifier{}
	orderService1 := NewOrderService(consoleNotifier)
	fmt.Println("Using Console Notifier:")
	orderService1.PlaceOrder()

	fmt.Println("\n--- Scenario 2 ---")
	// --- Scenario 2 ---
	// TODO: Create an EmailNotifier and inject it into a new OrderService
	emailNotifier := EmailNotifier{}
	orderService2 := NewOrderService(emailNotifier)
	fmt.Println("Using Email Notifier:")
	orderService2.PlaceOrder()
}
