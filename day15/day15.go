package main

import (
	"advent-2019/gridprint"
	"advent-2019/intcode"
	"advent-2019/point"
	"advent-2019/turtle"
	"fmt"
)

const (
	North = 1
	South = 2
	West  = 3
	East  = 4
)

const (
	HitWall       = 0
	MovedToEmpty  = 1
	MovedToOxygen = 2
)

const (
	UnexploredTile = 0
	EmptyTile      = 1
	WallTile       = 2
	OxygenTile     = 3
)

var directionToMoveCommand = map[turtle.Direction]int{
	turtle.Up:    North,
	turtle.Down:  South,
	turtle.Left:  West,
	turtle.Right: East,
}

var resultCodeToTileId = map[int]int{
	HitWall:       WallTile,
	MovedToOxygen: OxygenTile,
	MovedToEmpty:  EmptyTile,
}

type Grid map[point.Point]int

type System struct {
	tape   intcode.Tape
	turtle turtle.Turtle
	grid   Grid
}

func (system System) AdjacentPositions(pos point.Point) []point.Point {
	var positions []point.Point
	for _, dir := range turtle.Directions {
		posInDir := turtle.NextPosInDir(pos, dir)
		if system.CanMoveToPoint(posInDir) {
			positions = append(positions, posInDir)
		}
	}
	return positions
}

func (system *System) MoveRobot(direction turtle.Direction) bool {
	system.tape.ClearInput()
	system.tape.Input(directionToMoveCommand[direction])

	targetPos := turtle.NextPosInDir(system.turtle.Pos(), direction)

	result := system.tape.RunUntilNextOutput()

	system.grid[targetPos] = resultCodeToTileId[result]

	moveSuccess := result != HitWall
	if moveSuccess {
		system.turtle.MoveInDir(direction)
	}
	return moveSuccess
}

func (system System) CurrentTile() int {
	return system.grid[system.turtle.Pos()]
}

func (system *System) Explore(dir turtle.Direction, visited map[point.Point]bool) {
	result := system.MoveRobot(dir)
	if !result {
		return
	}

	defer system.MoveRobot(dir.Opposite())

	for _, nextDir := range turtle.Directions {
		posInDir := turtle.NextPosInDir(system.turtle.Pos(), nextDir)
		if visited[posInDir] || system.grid[posInDir] != UnexploredTile {
			continue
		}

		system.Explore(nextDir, visited)
	}
}

func (system System) CanMoveToPoint(pos point.Point) bool {
	tile := system.grid[pos]
	return tile != UnexploredTile && tile != WallTile
}

func pathLength(parents map[point.Point]point.Point, source point.Point, dest point.Point) int {
	fmt.Println("Finding path length")
	length := 0
	node := dest
	for {
		if node == source {
			return length
		}

		parent, ok := parents[node]
		if !ok {
			return -1
		}

		node = parent
		length++
	}
}

func (system System) DistanceBetween(source point.Point, dest point.Point) int {
	if !system.CanMoveToPoint(source) || !system.CanMoveToPoint(dest) {
		return -1
	}

	visited := map[point.Point]bool{}
	parents := map[point.Point]point.Point{}
	queue := point.Queue{}
	queue.Push(source)
	for !queue.IsEmpty() {
		pos := queue.Pop()

		fmt.Println("Visiting", pos)

		if pos == dest {
			return pathLength(parents, source, dest)
		}

		for _, dir := range turtle.Directions {
			posInDir := turtle.NextPosInDir(pos, dir)

			if !system.CanMoveToPoint(posInDir) || visited[posInDir] {
				continue
			}

			visited[posInDir] = true
			parents[posInDir] = pos
			queue.Push(posInDir)
		}
	}

	fmt.Println("Path not found")

	return -1
}

func (system System) GetOxygenLocation() point.Point {
	for pos, tile := range system.grid {
		if tile == OxygenTile {
			return pos
		}
	}
	return point.Point{-1, -1}
}

func (system System) Print() {
	gridprint.PrintGrid(system.grid, func(tile int, pos point.Point) {
		if pos == system.turtle.Pos() {
			fmt.Print("D")
		}

		switch tile {
		case UnexploredTile:
			fmt.Print(" ")
		case WallTile:
			fmt.Print("#")
		case EmptyTile:
			fmt.Print(".")
		case OxygenTile:
			fmt.Print("o")
		default:
			fmt.Print("?")
		}
	})
}

func part1() {
	system := System{tape: intcode.CreateBlankTape("advent-2019/day15.txt"), grid: map[point.Point]int{}, turtle: turtle.Turtle{}}
	visited := map[point.Point]bool{}
	for _, dir := range turtle.Directions {
		system.Explore(dir, visited)
	}
	oxygenLocation := system.GetOxygenLocation()
	fmt.Println(system.DistanceBetween(point.Point{0, 0}, oxygenLocation))
}

func part2() {
	system := System{tape: intcode.CreateBlankTape("advent-2019/day15.txt"), grid: map[point.Point]int{}, turtle: turtle.Turtle{}}
	visited := map[point.Point]bool{}
	for _, dir := range turtle.Directions {
		system.Explore(dir, visited)
	}
	oxygenLocation := system.GetOxygenLocation()
	queue := point.Queue{}
	queue.Push(oxygenLocation)
	minute := 0
	for !queue.IsEmpty() {
		fmt.Println("Minute", minute)
		system.Print()
		fmt.Println()
		currentQueue := queue
		queue = point.Queue{}
		for !currentQueue.IsEmpty() {
			pos := currentQueue.Pop()
			system.grid[pos] = OxygenTile
			for _, dir := range turtle.Directions {
				posInDir := turtle.NextPosInDir(pos, dir)
				if system.grid[posInDir] == EmptyTile {
					queue.Push(posInDir)
				}
			}
		}
		minute++
	}
	fmt.Println(minute)
}

func main() {
	part2()
}
