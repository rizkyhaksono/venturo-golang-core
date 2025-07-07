package main

import (
	"fmt"
	"time"
)

func produceMessage(ch chan string) {
	time.Sleep(5 * time.Second)
	ch <- "Message from the producer!" // Send message into the channel.
}

func main() {
	messageChannel := make(chan string) // Create a channel of type string.

	go produceMessage(messageChannel)

	fmt.Println("Waiting for a message...")
	// Block and wait to receive a message from the channel.
	receivedMessage := <-messageChannel
	fmt.Printf("Received: %s\n", receivedMessage)
}
