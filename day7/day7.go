package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"util/datafile"
)

const (
	addOpcode         int = 1
	multiplyOpcode    int = 2
	inputOpcode       int = 3
	outputOpcode      int = 4
	jumpIfTrueOpcode  int = 5
	jumpIfFalseOpcode int = 6
	lessThanOpcode    int = 7
	equalsOpcode      int = 8
	haltOpcode        int = 99
)

const (
	positionMode  = 0
	immediateMode = 1
)

type param struct {
	value int
	mode  int
}

type tape struct {
	data   []int
	cursor int
	input  chan int
	output *chan int
}

func getChar(s string, pos int) byte {
	if pos >= len(s) || pos < 0 {
		return '0'
	}

	return s[len(s)-pos-1]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func decodeOpcode(value int) int {
	valueString := strconv.Itoa(value)

	opcodeString := string(getChar(valueString, 1)) + string(getChar(valueString, 0))
	opcode, err := strconv.Atoi(opcodeString)
	if err != nil {
		log.Fatal(err)
	}

	return opcode
}

func prompt() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func (t tape) IsHalted() bool {
	return t.Value() == haltOpcode
}

// GetParams extracts [paramCount] params, and advances the cursor to the instruction that will happen next
func (t *tape) GetParams(instructionValue int, paramCount int) []param {
	instructionString := strconv.Itoa(instructionValue)
	params := make([]param, paramCount)

	cursor := t.cursor + 1
	for i := 0; i < paramCount; i++ {
		paramValue := t.data[cursor+i]
		modeString := getChar(instructionString, i+2)
		mode, err := strconv.Atoi(string(modeString))
		if err != nil {
			log.Fatal(err)
		}
		params[i] = param{paramValue, mode}
	}

	t.cursor += paramCount + 1

	return params
}

// Value returns the value/opcode at the cursor
func (t tape) Value() int {
	return t.data[t.cursor]
}

// Resolve returns the value of a parameter, based on its mode
func (t tape) Resolve(p param) int {
	switch p.mode {
	case positionMode:
		return t.data[p.value]
	case immediateMode:
		return p.value
	default:
		log.Fatal("Invalid mode:", p.mode)
		return 0 // this will never be hit
	}
}

// Output returns the value at the first index, aka the output.
// This return value is invalid if the tape has not been run
func (t tape) Output() int {
	return t.data[0]
}

// Run will run the tape from the current data/cursor until it halts or hits an unknown opcode
func (t *tape) Run(wg *sync.WaitGroup, i int) {
	if wg != nil {
		wg.Add(1)
		defer wg.Done()
	}

	for !t.IsHalted() {
		value := t.Value()
		opcode := decodeOpcode(value)

		// fmt.Println("Processing value:", value)
		// fmt.Println("Current data:", t.data)
		// fmt.Println("Opcode:", opcode)

		switch opcode {
		case addOpcode:
			{
				params := t.GetParams(value, 3)
				a := t.Resolve(params[0])
				b := t.Resolve(params[1])
				outputIndex := params[2].value
				t.data[outputIndex] = a + b
			}
		case multiplyOpcode:
			{
				params := t.GetParams(value, 3)
				a := t.Resolve(params[0])
				b := t.Resolve(params[1])
				outputIndex := params[2].value
				t.data[outputIndex] = a * b
			}
		case inputOpcode:
			{
				params := t.GetParams(value, 1)
				destination := params[0].value
				//fmt.Println(i, "Waiting for input")
				t.data[destination] = <-t.input
				//fmt.Println(i, "Got input:", t.data[destination])
			}
		case outputOpcode:
			{
				params := t.GetParams(value, 1)
				value := t.Resolve(params[0])
				//fmt.Println(i, "Sending output", value)
				*t.output <- value
				//fmt.Println(i, "Sent output")
			}
		case jumpIfTrueOpcode:
			{
				params := t.GetParams(value, 2)
				testValue := t.Resolve(params[0])
				jumpIndex := t.Resolve(params[1])
				if testValue != 0 {
					t.cursor = jumpIndex
				}
			}
		case jumpIfFalseOpcode:
			{
				params := t.GetParams(value, 2)
				testValue := t.Resolve(params[0])
				jumpIndex := t.Resolve(params[1])
				if testValue == 0 {
					t.cursor = jumpIndex
				}
			}
		case lessThanOpcode:
			{
				params := t.GetParams(value, 3)
				a := t.Resolve(params[0])
				b := t.Resolve(params[1])
				destination := params[2].value

				writeValue := 0
				if a < b {
					writeValue = 1
				}

				t.data[destination] = writeValue
			}
		case equalsOpcode:
			{
				params := t.GetParams(value, 3)
				a := t.Resolve(params[0])
				b := t.Resolve(params[1])
				destination := params[2].value

				writeValue := 0
				if a == b {
					writeValue = 1
				}
				t.data[destination] = writeValue
			}
		default:
			log.Fatal("Invalid opcode: ", opcode)
		}
	}
	//fmt.Println(i, "Done")
}

func getTapeData() []int {
	file := datafile.Open("advent-2019/day7.txt")
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
	return tape{data, 0, make(chan int), nil}
}

func createTapeCopy(data []int) tape {
	dataCopy := make([]int, len(data))
	copy(dataCopy, data)
	return tape{dataCopy, 0, make(chan int), nil}
}

func permutations(items []int, callback func([]int), i int) {
	if i > len(items) {
		callback(items)
		return
	}

	permutations(items, callback, i+1)

	for j := i + 1; j < len(items); j++ {
		items[i], items[j] = items[j], items[i]
		permutations(items, callback, i+1)
		items[i], items[j] = items[j], items[i]
	}
}

func wrap(v, l int) int {
	for v < 0 {
		v += l
	}

	return v % l
}

func linkAmps(amplifiers []tape) {
	for i := 0; i < len(amplifiers); i++ {
		amp := &amplifiers[i]
		nextAmp := &amplifiers[wrap(i+1, len(amplifiers))]
		amp.output = &nextAmp.input
	}
}

func startAmps(amplifiers []tape, wg *sync.WaitGroup) {
	lastIndex := len(amplifiers)-1

	for i := 0; i < len(amplifiers); i++ {
		amp := &amplifiers[i]
		if i < lastIndex {
			go amp.Run(wg, i)
		} else {
			go amp.Run(nil, i)
		}
	}

}

func part1() {
	data := getTapeData()
	highestOutput := -1
	var highestOutputPhases []int
	permutations([]int{0, 1, 2, 3, 4}, func(phases []int) {
		amplifiers := []tape{
			createTapeCopy(data),
			createTapeCopy(data),
			createTapeCopy(data),
			createTapeCopy(data),
			createTapeCopy(data),
		}
		linkAmps(amplifiers)

		var wg sync.WaitGroup
		startAmps(amplifiers, &wg)

		for i, phase := range phases {
			amplifiers[i].input <- phase
		}

		amplifiers[0].input <- 0

		wg.Wait()
		output := <-amplifiers[0].input

		if highestOutput == -1 || output > highestOutput {
			highestOutput = output
			highestOutputPhases = phases
		}
	}, 0)
	fmt.Println("Highest output:", highestOutput)
	fmt.Println("Phases:", highestOutputPhases)
}

func part2() {
	data := getTapeData()
	highestOutput := -1
	permutations([]int{5, 6, 7, 8, 9}, func(phases []int) {
		amplifiers := []tape{
			createTapeCopy(data),
			createTapeCopy(data),
			createTapeCopy(data),
			createTapeCopy(data),
			createTapeCopy(data),
		}
		linkAmps(amplifiers)

		var wg sync.WaitGroup
		startAmps(amplifiers, &wg)

		for i, phase := range phases {
			amplifiers[i].input <- phase
		}

		amplifiers[0].input <- 0

		wg.Wait()
		output := <-amplifiers[0].input

		if highestOutput == -1 || output > highestOutput {
			highestOutput = output
		}
	}, 0)
	fmt.Println(highestOutput)
}

func main() {
	part2()
}
