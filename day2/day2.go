package main

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"
	"util/datafile"
)

const (
	add      int = 1
	multiply int = 2
	halt     int = 99
)

type tape struct {
	data   []int
	cursor int
}

// GetParams extracts [paramCount] params, and advances the cursor to the instruction that will happen next
func (t *tape) GetParams(paramCount int) []int {
	params := make([]int, paramCount)

	cursor := (*t).cursor + 1
	for i := 0; i < paramCount; i++ {
		params[i] = (*t).data[cursor+i]
	}

	(*t).cursor += (paramCount + 1)

	return params
}

// Value returns the value/opcode at the cursor
func (t tape) Value() int {
	return t.data[t.cursor]
}

// Output returns the value at the first index, aka the output.
// This return value is invalid if the tape has not been run
func (t tape) Output() int {
	return t.data[0]
}

// Run will run the tape from the current data/cursor until it halts or hits an unknown opcode
func (t tape) Run() {
	for t.Value() != halt {
		opcode := t.Value()

		switch opcode {
		case add:
			{
				params := t.GetParams(3)
				inputIndexA := params[0]
				inputIndexB := params[1]
				outputIndex := params[2]
				t.data[outputIndex] = t.data[inputIndexA] + t.data[inputIndexB]
			}
		case multiply:
			{
				params := t.GetParams(3)
				inputIndexA := params[0]
				inputIndexB := params[1]
				outputIndex := params[2]
				t.data[outputIndex] = t.data[inputIndexA] * t.data[inputIndexB]
			}
		default:
			log.Fatal("Invalid opcode: ", opcode)
		}
	}
}

func getTapeData() []int {
	file := datafile.Open("advent-2019/day2.txt")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()

	numberStrings := strings.Split(line, ",")
	data := make([]int, len(numberStrings))
	for i, numberString := range numberStrings {
		number, err := strconv.Atoi(numberString)
		if err != nil {
			log.Fatal(err)
		}
		data[i] = int(number)
	}

	return data
}

func createBlankTape() tape {
	data := getTapeData()
	return tape{data, 0}
}

func part1() {
	fmt.Println("Creating tape")

	t := createBlankTape()

	t.data[1] = 12
	t.data[2] = 2

	fmt.Println("Running tape")

	t.Run()

	fmt.Println("Halt")
	fmt.Println(t.data)
}

func findNounAndVerb() (int, int) {
	baseData := getTapeData()

	desiredOutput := 19690720

	for noun := 0; noun < 100; noun++ {
		for verb := 0; verb < 100; verb++ {
			data := make([]int, len(baseData))
			copy(data, baseData)

			data[1] = noun
			data[2] = verb

			t := tape{data, 0}
			t.Run()

			if t.Output() == desiredOutput {
				return noun, verb
			}
		}
	}

	return -1, -1
}

func part2() {
	noun, verb := findNounAndVerb()

	if noun == -1 || verb == -1 {
		log.Fatal("Could not find valid combination")
	}

	fmt.Println((100 * noun) + verb)
}

func main() {
	part2()
}
