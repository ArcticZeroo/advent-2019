package gridprint

import (
	"advent-2019/point"
	"fmt"
)

func PrintGrid(grid map[point.Point]int, printItem func(value int, pos point.Point)) {
	smallestX, smallestY, largestX, largestY := -1, -1, -1, -1
	for p := range grid {
		if smallestX == -1 || p.X < smallestX {
			smallestX = p.X
		}
		if smallestY == -1 || p.Y < smallestY {
			smallestY = p.Y
		}
		if largestX == -1 || p.X > largestX {
			largestX = p.X
		}
		if largestY == -1 || p.Y > largestY {
			largestY = p.Y
		}
	}

	for y := largestY; y >= smallestY; y-- {
		for x := smallestX; x <= largestX; x++ {
			pos := point.Point{X: x, Y: y}
			value := grid[pos]
			printItem(value, pos)
		}
		fmt.Println()
	}
}
