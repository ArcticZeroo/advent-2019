package main

import (
	"advent-2019/smath"
	"bufio"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	"util/datafile"
)

const (
	Ore  = "ORE"
	Fuel = "FUEL"
)

var recipeRegexp = regexp.MustCompile(`(.+) => (.+)`)

type Component struct {
	id    string
	count int
}

type Recipe struct {
	inputs []Component
	output Component
}

type Cookbook map[string]Recipe

type Kitchen struct {
	cookbook            Cookbook
	inventory           Inventory
	usedIngredients     Inventory
	producedIngredients Inventory
}

func (k Kitchen) UsedOre() int {
	return k.usedIngredients[Ore]
}

type Inventory map[string]int

func (i Inventory) Needed(component Component) int {
	return smath.MaxInt(0, component.count-i[component.id])
}

func (i Inventory) Consume(component Component) {
	i[component.id] -= component.count
}

func (i Inventory) Produce(component Component) {
	i[component.id] += component.count
}

func (k Kitchen) UseIngredient(ingredient Component) {
	k.usedIngredients[ingredient.id] += ingredient.count
	k.inventory.Consume(ingredient)
}

func (k Kitchen) ProduceRecipeOutput(ingredient Component) {
	k.inventory.Produce(ingredient)
	k.producedIngredients.Produce(ingredient)
}

func (k *Kitchen) ProduceIngredient(output Component) {
	if output.id == Ore {
		return
	}

	needed := k.inventory.Needed(output)
	recipe := k.cookbook[output.id]
	recipeCookCount := int(math.Ceil(float64(needed) / float64(recipe.output.count)))
	for _, input := range recipe.inputs {
		ingredient := Component{input.id, input.count * recipeCookCount}
		k.ProduceIngredient(ingredient)
		k.UseIngredient(ingredient)
	}
	k.ProduceRecipeOutput(Component{output.id, recipe.output.count * recipeCookCount})
}

func componentFromString(s string) Component {
	pieces := strings.Split(s, " ")
	count, err := strconv.Atoi(pieces[0])
	if err != nil {
		log.Fatal(err)
	}
	return Component{pieces[1], count}
}

func recipeFromString(s string) Recipe {
	match := recipeRegexp.FindStringSubmatch(s)
	inputStrings := strings.Split(match[1], ", ")

	inputs := make([]Component, len(inputStrings))
	output := componentFromString(match[2])

	for i, s := range inputStrings {
		inputs[i] = componentFromString(s)
	}

	return Recipe{inputs, output}
}

func getCookbookFromFile() Cookbook {
	cookbook := Cookbook{}

	file := datafile.Open("advent-2019/day14.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		recipe := recipeFromString(line)
		cookbook[recipe.output.id] = recipe
	}

	return cookbook
}

func fuelCost(cookbook Cookbook, n int) int {
	inventory := Inventory{}
	usedIngredients := Inventory{}
	producedIngredients := Inventory{}
	fuel := Component{Fuel, n}
	kitchen := Kitchen{cookbook, inventory, usedIngredients, producedIngredients}
	kitchen.ProduceIngredient(fuel)
	return kitchen.UsedOre()
}

func findFuelCount(cookbook Cookbook, lowerCount, upperCount int) int {
	if lowerCount >= upperCount {
		return lowerCount
	}

	midCount := lowerCount + ((upperCount - lowerCount) / 2)

	cost := fuelCost(cookbook, midCount + 1)
	target := int(1e12)

	if cost == target {
		return midCount
	}

	if cost > target {
		return findFuelCount(cookbook, lowerCount, midCount)
	}

	return findFuelCount(cookbook, midCount + 1, upperCount)
}

func solve() {
	cookbook := getCookbookFromFile()
	singleFuelCost := fuelCost(cookbook, 1)
	targetOre := int(1e12)
	fuelCount := targetOre / singleFuelCost
	fmt.Println("Base fuel producedCount:", fuelCount)
	fmt.Println("Max fuel count:", findFuelCount(cookbook, fuelCount, fuelCount * 2))
}

func main() {
	solve()
}
