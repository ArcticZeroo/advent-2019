package intcode

import (
	"bufio"
	"intqueue"
	"log"
	"strconv"
	"strings"
	"util/datafile"
)

const (
	addOpcode            int = 1
	multiplyOpcode       int = 2
	inputOpcode          int = 3
	outputOpcode         int = 4
	jumpIfTrueOpcode     int = 5
	jumpIfFalseOpcode    int = 6
	lessThanOpcode       int = 7
	equalsOpcode         int = 8
	relativeAdjustOpcode int = 9
	haltOpcode           int = 99
)

const (
	positionMode  = 0
	immediateMode = 1
	relativeMode  = 2
)

type param struct {
	value int
	mode  int
}

type Tape struct {
	data         []int
	cursor       int
	input        intqueue.Queue
	output       intqueue.Queue
	relativeBase int
}

func getChar(s string, pos int) byte {
	if pos >= len(s) || pos < 0 {
		return '0'
	}

	return s[len(s)-pos-1]
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

func (t Tape) IsHalted() bool {
	return t.Value() == haltOpcode
}

// GetParams extracts [paramCount] params, and advances the cursor to the instruction that will happen next
func (t *Tape) GetParams(instructionValue int, paramCount int) []param {
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
func (t Tape) Value() int {
	return t.data[t.cursor]
}

// Resolve returns the address to a value. This value may be used as a source or destination
func (t *Tape) Resolve(p param) *int {
	switch p.mode {
	case positionMode:
		return &t.data[p.value]
	case immediateMode:
		return &p.value
	case relativeMode:
		return &t.data[p.value+t.relativeBase]
	default:
		log.Fatal("Invalid mode:", p.mode)
		return nil // this will never be hit
	}
}

// First returns the value at the first index, aka the output.
// This return value is invalid if the tape has not been run
func (t Tape) First() int {
	return t.data[0]
}

func (t *Tape) Input(x int) {
	t.input.Push(x)
}

func (t Tape) Output() intqueue.Queue {
	return t.output
}

func (t *Tape) RunNextInstruction() {
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
			dest := t.Resolve(params[2])
			*dest = *a + *b
		}
	case multiplyOpcode:
		{
			params := t.GetParams(value, 3)
			a := t.Resolve(params[0])
			b := t.Resolve(params[1])
			dest := t.Resolve(params[2])
			*dest = (*a) * (*b)
		}
	case inputOpcode:
		{
			params := t.GetParams(value, 1)
			dest := t.Resolve(params[0])
			*dest = t.input.Pop()
		}
	case outputOpcode:
		{
			params := t.GetParams(value, 1)
			value := t.Resolve(params[0])
			t.output.Push(*value)
		}
	case jumpIfTrueOpcode:
		{
			params := t.GetParams(value, 2)
			testValue := t.Resolve(params[0])
			jumpIndex := t.Resolve(params[1])
			if *testValue != 0 {
				t.cursor = *jumpIndex
			}
		}
	case jumpIfFalseOpcode:
		{
			params := t.GetParams(value, 2)
			testValue := t.Resolve(params[0])
			jumpIndex := t.Resolve(params[1])
			if *testValue == 0 {
				t.cursor = *jumpIndex
			}
		}
	case lessThanOpcode:
		{
			params := t.GetParams(value, 3)
			a := t.Resolve(params[0])
			b := t.Resolve(params[1])
			destination := t.Resolve(params[2])

			writeValue := 0
			if *a < *b {
				writeValue = 1
			}

			*destination = writeValue
		}
	case equalsOpcode:
		{
			params := t.GetParams(value, 3)
			a := t.Resolve(params[0])
			b := t.Resolve(params[1])
			destination := t.Resolve(params[2])

			writeValue := 0
			if *a == *b {
				writeValue = 1
			}
			*destination = writeValue
		}
	case relativeAdjustOpcode:
		{
			params := t.GetParams(value, 1)
			increment := t.Resolve(params[0])
			t.relativeBase += *increment
		}
	default:
		log.Fatal("Invalid opcode: ", opcode)
	}
}

// Run will run the tape from the current data/cursor until it halts or hits an unknown opcode
func (t *Tape) RunUntilHalt() {
	for !t.IsHalted() {
		t.RunNextInstruction()
	}
}

func (t *Tape) RunUntilNextOutput() int {
	for !t.IsHalted() && t.output.Empty() {
		t.RunNextInstruction()
	}

	if t.output.Empty() {
		return 0
	}

	return t.output.Pop()
}

// GetTapeData returns data for a tape from the given path
func GetTapeData(path string) []int {
	file := datafile.Open(path)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()

	numberStrings := strings.Split(line, ",")
	data := make([]int, len(numberStrings)*8)
	for i, numberString := range numberStrings {
		number, err := strconv.Atoi(numberString)
		if err != nil {
			log.Fatal(err)
		}
		data[i] = int(number)
	}

	return data
}

// CreateBlankTape returns a blank tape based on the given input
func CreateBlankTape(path string) Tape {
	data := GetTapeData(path)
	return Tape{data, 0, intqueue.Queue{}, intqueue.Queue{}, 0}
}

// CreateTapeCopy creates a new tape with the given data copied
func CreateTapeCopy(data []int) Tape {
	dataCopy := make([]int, len(data))
	copy(dataCopy, data)
	return Tape{dataCopy, 0, intqueue.Queue{}, intqueue.Queue{}, 0}
}
