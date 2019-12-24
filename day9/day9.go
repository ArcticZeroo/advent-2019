package main

import (
	"advent-2019/intcode"
	"fmt"
)

func part1() {
	tape := intcode.CreateBlankTape("advent-2019/day9.txt")
	tape.Input(2)
	tape.RunUntilHalt()
	output := tape.Output()
	for !output.Empty() {
		fmt.Println(output.Pop())
	}
}

func main() {
	part1()
}
