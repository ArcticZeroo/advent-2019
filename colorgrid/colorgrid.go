package colorgrid

import (
	"advent-2019/point"
	"fmt"
	"github.com/logrusorgru/aurora"
)

type Color int

const (
	Black = 0
	White = 1
)

type Grid map[point.Point]Color

func (g Grid) Print() {
	smallestX, smallestY, largestX, largestY := -1, -1, -1, -1
	for p := range g {
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
			color := g[point.Point{X: x, Y: y}]
			switch color {
			case Black:
				fmt.Print(aurora.White(" ").BgBlack())
			case White:
				fmt.Print(aurora.Black(" ").BgWhite())
			default:
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}
