package main

import (
	"fmt"
	"strconv"
)

func meetsRequirements(x int) bool {
	xs := strconv.Itoa(x)
	hasSeenAdjacent := false
	adjacencyLength := 0
	last := -1
	for _, digitRune := range xs {
		digit, _ := strconv.Atoi(string(digitRune))

		if digit < last {
			return false
		}

		if last == digit {
			adjacencyLength++
		} else {
			if adjacencyLength == 2 {
				hasSeenAdjacent = true
			}

			adjacencyLength = 1
		}

		last = digit
	}
	return hasSeenAdjacent || adjacencyLength == 2
}

func main() {
	fmt.Println(meetsRequirements(112233))
	fmt.Println(meetsRequirements(123444))
	fmt.Println(meetsRequirements(111122))

	rangeStart := 271973
	rangeEnd := 785961

	total := 0
	for i := rangeStart; i <= rangeEnd; i++ {
		if meetsRequirements(i) {
			total++
		}
	}
	fmt.Println(total)
}
