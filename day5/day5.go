package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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

	// paramModes := make([]int, max(0, len(valueString)-2))
	// for i := 2; i < len(valueString); i++ {
	// 	char := getChar(valueString, i)
	// 	digit, err := strconv.Atoi(string(char))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	paramModes[i-2] = digit
	// }

	// return instruction{opcode, paramModes}
}

func prompt() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
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

	t.cursor += (paramCount + 1)

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
func (t tape) Run() {
	for t.Value() != haltOpcode {
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

				var value int
				for {
					fmt.Print("Instruction input: ")
					input := prompt()
					possibleValue, err := strconv.Atoi(input)
					if err == nil {
						value = possibleValue
						break
					}
					fmt.Println("Invalid input")
				}

				t.data[destination] = value
			}
		case outputOpcode:
			{
				params := t.GetParams(value, 1)
				value := t.Resolve(params[0])
				fmt.Println(value)
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
}

func getTapeData() []int {
	file := datafile.Open("advent-2019/day5.txt")
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

func main() {
	// fmt.Println("Creating tape")

	t := createBlankTape()

	// fmt.Println("Running tape")

	t.Run()

	// fmt.Println("Halt")
	// fmt.Println(t.data)
}
