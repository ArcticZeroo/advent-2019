package main

import (
	"advent-2019/colorgrid"
	"advent-2019/intcode"
	"advent-2019/point"
	"advent-2019/turtle"
)

const (
	turnLeft = 0
	turnRight = 1
)

const (
	paintBlack = 0
	paintWhite = 1
)

func solve() {
	grid := colorgrid.Grid{}
	grid[point.Point{}] = colorgrid.White
	robot := turtle.Turtle{}
	tape := intcode.CreateBlankTape("advent-2019/day11.txt")
	for !tape.IsHalted() {
		tape.Input(int(grid[robot.Pos()]))
		color := tape.RunUntilNextOutput()
		turnDir := tape.RunUntilNextOutput()

		if tape.IsHalted() {
			break
		}

		grid[robot.Pos()] = colorgrid.Color(color)

		if turnDir == turnLeft {
			robot.TurnLeft()
		} else {
			robot.TurnRight()
		}

		robot.MoveForward()
	}
	grid.Print()
}

func main() {
	solve()
}