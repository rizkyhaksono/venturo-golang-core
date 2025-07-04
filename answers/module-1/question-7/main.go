package main

import "fmt"

type Book struct {
	Title  string
	Author string
	Pages  int16
	ISBN   string
}

func main() {
	// TODO: Create a slice of Book structs and add at least 3 books.
	library := []Book{
		{Title: "Buku pertama", Author: "Penulis A", Pages: 100, ISBN: "123-456-789"},
		{Title: "Buku kedua", Author: "Penulis B", Pages: 200, ISBN: "987-654-321"},
		{Title: "Buku ketiga", Author: "Penulis C", Pages: 300, ISBN: "456-789-123"},
	}

	fmt.Println("Library Titles:")
	// TODO: Write a for loop to iterate over the library
	// and print the Title of each book.
	for _, book := range library {
		fmt.Printf("Title: %s, Author %s, Pages: %d, ISBN: %s\n", book.Title, book.Author, book.Pages, book.ISBN)
	}
}
