package main

import (
	"fmt"
	"log"
	"strconv"
	"util/datafile"

	. "github.com/logrusorgru/aurora"
)

type imageLayer []int

const (
	width          = 25
	height         = 6
	pixelsPerLayer = width * height
	black          = 0
	white          = 1
	transparent    = 2
)

func (layer imageLayer) Print() {
	for y := 0; y < height; y++ {
		offsetX := y * width
		for x := 0; x < width; x++ {
			pixel := layer[offsetX+x]
			switch pixel {
			case black:
				fmt.Print(White(" ").BgBlack())
			case white:
				fmt.Print(Black(" ").BgWhite())
			case transparent:
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func getLayerIndex(i int) int {
	return i / pixelsPerLayer
}

func getPixelIndex(i int) int {
	return i % pixelsPerLayer
}

func getLayersFromInput() []imageLayer {
	line := datafile.OpenFirstLine("advent-2019/day8.txt")
	layerCount := len(line) / pixelsPerLayer
	image := make([]imageLayer, layerCount)

	for i, c := range line {
		value, err := strconv.Atoi(string(c))
		if err != nil {
			log.Fatal(err)
		}
		layerIndex := getLayerIndex(i)
		if image[layerIndex] == nil {
			image[layerIndex] = make(imageLayer, pixelsPerLayer)
		}
		pixelIndex := getPixelIndex(i)
		image[layerIndex][pixelIndex] = value
	}

	return image
}

func part1() {
	image := getLayersFromInput()

	fewestZeroCount := -1
	var fewestZeroCounter map[int]int
	for _, layer := range image {
		numberCounts := map[int]int{}
		for _, x := range layer {
			numberCounts[x]++
		}
		zeroCount := numberCounts[0]
		if fewestZeroCount == -1 || zeroCount < fewestZeroCount {
			fewestZeroCount = zeroCount
			fewestZeroCounter = numberCounts
		}
	}

	if fewestZeroCount == -1 || fewestZeroCounter == nil {
		log.Fatal("Layer with fewest zeros not found")
	}

	fmt.Println(fewestZeroCounter[1] * fewestZeroCounter[2])
}

func getRenderedPixel(layers []imageLayer, pos int) int {
	for i := 0; i < len(layers); i++ {
		pixel := layers[i][pos]
		if pixel != transparent {
			return pixel
		}
	}
	return transparent
}

func part2() {
	layers := getLayersFromInput()
	renderedImage := make(imageLayer, pixelsPerLayer)
	for i := 0; i < pixelsPerLayer; i++ {
		renderedImage[i] = getRenderedPixel(layers, i)
	}
	renderedImage.Print()
}

func main() {
	part1()
	part2()
}
