package main

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"
	"util/datafile"
)

type direction int
type overlapGrid map[point]int

const (
	up    direction = 0
	down  direction = 1
	left  direction = 2
	right direction = 3
)

var charToDirection = map[byte]direction{
	'U': up,
	'D': down,
	'L': left,
	'R': right,
}

type lineInstruction struct {
	dir    direction
	length int
}

type point struct {
	x, y int
}

type line struct {
	instructions []lineInstruction
	position     point
	steps        int
}

func (l *line) Step(dir direction) {
	position := (*l).position

	toIncrement := &position.x
	if dir == up || dir == down {
		toIncrement = &position.y
	}

	stepSize := 1
	if dir == down || dir == left {
		stepSize = -1
	}

	*toIncrement += stepSize

	(*l).position = position
	(*l).steps++
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (l *line) MoveNextAndFindShortestIntersection(grid overlapGrid) int {
	/*
		shortestDistance := -1
		for i := 0; i < instruction.length; i++ {
			*toIncrement += (1 * multiplier)
			if grid[position.x][position.y] {
				intersectionDistance := absInt(position.x) + absInt(position.y)
				if shortestDistance == -1 || intersectionDistance < shortestDistance {
					shortestDistance = intersectionDistance
				}
			}
		}

	return shortestDistance*/
	return -1
}

func getScannerLine(scanner *bufio.Scanner) string {
	scanner.Scan()
	return scanner.Text()
}

func getLineInstructions(line string) []lineInstruction {
	instructionStrings := strings.Split(line, ",")
	instructions := make([]lineInstruction, len(instructionStrings))
	for i, instructionString := range instructionStrings {
		directionString := instructionString[0]
		dir := charToDirection[directionString]
		lengthString := instructionString[1:]
		length, err := strconv.Atoi(lengthString)
		if err != nil {
			log.Fatal(err)
		}
		instructions[i] = lineInstruction{dir, length}
	}
	return instructions
}

func getNextLine(scanner *bufio.Scanner) line {
	instructions := getLineInstructions(getScannerLine(scanner))
	fmt.Println(instructions)
	return line{instructions, point{}, 0}
}

func printLineGrid(grid overlapGrid) {
	smallestX, smallestY, largestX, largestY := -1, -1, -1, -1
	for point := range grid {
		if smallestX == -1 || point.x < smallestX {
			smallestX = point.x
		}
		if smallestY == -1 || point.y < smallestY {
			smallestY = point.y
		}
		if largestX == -1 || point.x > largestX {
			largestX = point.x
		}
		if largestY == -1 || point.y > largestY {
			largestY = point.y
		}
	}

	for y := largestY; y >= 0; y-- {
		for x := 0; x <= largestX; x++ {
			if (grid[point{x, y}] > 0) {
				fmt.Print("x")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func part1() {
	file := datafile.Open("advent-2019/day3.txt")

	grid := overlapGrid{}

	scanner := bufio.NewScanner(file)
	lineA := getNextLine(scanner)
	lineB := getNextLine(scanner)

	for _, instruction := range lineA.instructions {
		for i := 0; i < instruction.length; i++ {
			lineA.Step(instruction.dir)
			grid[lineA.position] = lineA.steps
		}
	}

	shortestDistance := -1
	for _, instruction := range lineB.instructions {
		for i := 0; i < instruction.length; i++ {
			lineB.Step(instruction.dir)
			if grid[lineB.position] > 0 {
				distance := grid[lineB.position] + lineB.steps
				if shortestDistance == -1 || distance < shortestDistance {
					shortestDistance = distance
				}
			}
		}
	}

	fmt.Println(shortestDistance)
}

func main() {
	part1()
}
