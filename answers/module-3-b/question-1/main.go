package main

import "fmt"

type Speaker interface {
	Speak() string
}

type Human struct{}
type Dog struct{}

func (h Human) Speak() string {
	return "Hello"
}

func (d Dog) Speak() string {
	return "Woof"
}

// This function can accept any type that satisfies the Speaker interface.
func makeItSpeak(s Speaker) {
	fmt.Println(s.Speak())
}

func main() {
	human := Human{}
	dog := Dog{}

	makeItSpeak(human) // Expected: Hello
	makeItSpeak(dog)   // Expected: Woof
}
