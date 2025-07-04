package main

import "fmt"

type Month uint8

// TODO: Use iota to define the months from January to December.
const (
	_ Month = iota // Ignore 0
	January
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
)

// TODO: Implement this function to convert a Month to its string representation.
func GetMonthName(m Month) string {
	// Hint: A switch statement or a slice of strings would work well here.
	switch m {
	case January:
		return "January"
	case February:
		return "February"
	case March:
		return "March"
	case April:
		return "April"
	case May:
		return "May"
	case June:
		return "June"
	case July:
		return "July"
	case August:
		return "August"
	case September:
		return "September"
	case October:
		return "October"
	case November:
		return "November"
	case December:
		return "December"
	default:
		return "Unknown Month"
	}
}

func main() {
	fmt.Println(GetMonthName(January)) // Expected: January
}
