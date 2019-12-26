package main

import (
	"advent-2019/intcode"
	"advent-2019/point"
	"advent-2019/smath"
	"fmt"
	"github.com/logrusorgru/aurora"
)

const (
	emptyTile  = 0
	wallTile   = 1
	blockTile  = 2
	paddleTile = 3
	ballTile   = 4
)

type Game struct {
	grid        map[point.Point]int
	tape        intcode.Tape
	ballVel     point.Point
	ballPos     point.Point
	paddlePos   point.Point
	score       int
}

func (g *Game) setTile(pos point.Point, tile int) {
	switch tile {
	case ballTile:
		{
			if g.paddlePos.X != -1 {
				dir := pos.Subtract(g.paddlePos)
				dir.X /= smath.AbsInt(dir.X)
				g.tape.ClearInput()
				g.tape.Input(dir.X)
			}
		}
	case paddleTile:
		{
			g.paddlePos = pos
		}
	}

	g.grid[pos] = tile
}

func (g *Game) Play() {
	tape := &g.tape
	tape.Set(0, 2)

	for !tape.IsHalted() {
		tape.Input(0)

		x := tape.RunUntilNextOutput()
		y := tape.RunUntilNextOutput()
		tileOrScore := tape.RunUntilNextOutput()

		if x == -1 && y == 0 {
			g.score = tileOrScore
			continue
		}

		pos := point.Point{x, y}
		g.setTile(pos, tileOrScore)
	}
}

func part1() {
	tape := intcode.CreateBlankTape("advent-2019/day13.txt")
	grid := map[point.Point]int{}
	for !tape.IsHalted() {
		x := tape.RunUntilNextOutput()
		y := tape.RunUntilNextOutput()
		tile := tape.RunUntilNextOutput()
		grid[point.Point{x, y}] = tile
	}
	blockCount := 0
	for _, tile := range grid {
		if tile == blockTile {
			blockCount++
		}
	}
	fmt.Println(blockCount)
}

func print(grid map[point.Point]int) {
	largestX, largestY := -1, -1
	for p := range grid {
		if largestX == -1 || p.X > largestX {
			largestX = p.X
		}
		if largestY == -1 || p.Y > largestY {
			largestY = p.Y
		}
	}

	for y := 0; y < largestY; y++ {
		for x := 0; x < largestX; x++ {
			tile := grid[point.Point{x, y}]
			switch tile {
			case emptyTile:
				fmt.Print(" ")
			case wallTile:
				fmt.Print(aurora.White(" ").BgWhite())
			case blockTile:
				fmt.Print(aurora.BrightBlack(" ").BgBrightBlack())
			case paddleTile:
				fmt.Print("-")
			case ballTile:
				fmt.Print("o")
			}
		}
		fmt.Println()
	}
}

func part2() {
	tape := intcode.CreateBlankTape("advent-2019/day13.txt")
	game := Game{
		tape:      tape,
		grid:      map[point.Point]int{},
		ballVel:   point.Point{-1, -1},
		ballPos:   point.Point{-1, -1},
		paddlePos: point.Point{-1, -1},
		score:     0,
	}
	game.Play()
	fmt.Println("Score:", game.score)
}

func main() {
	part2()
}
