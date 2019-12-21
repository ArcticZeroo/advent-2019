package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func calculateFuel(mass int) int {
	return (mass / 3) - 2
}

func calculateTotalFuel(mass int) (total int) {
	fuel := calculateFuel(mass)
	for fuel > 0 {
		total += fuel
		fuel = calculateFuel(fuel)
	}
	return
}

func part1() {
	fmt.Println(os.Getenv("GOPATH"))

	path := filepath.Join(os.Getenv("GOPATH"), "./data/advent-2019/day1.txt")

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var total int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		value, err := strconv.Atoi(line)
		if err != nil {
			log.Fatal(err)
		}
		total += calculateFuel(value)
	}

	fmt.Println("Total:", total)
}

func part2() {
	fmt.Println(os.Getenv("GOPATH"))

	path := filepath.Join(os.Getenv("GOPATH"), "./data/advent-2019/day1.txt")

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var total int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		value, err := strconv.Atoi(line)
		if err != nil {
			log.Fatal(err)
		}
		total += calculateTotalFuel(value)
	}

	fmt.Println("Total:", total)
}

func main() {
	fmt.Println("Advent Day 1")

	fmt.Println("Part 1")
	part1()
	fmt.Println("Part 2")
	part2()
}
