package main

import (
	"bufio"
	"fmt"
	"math"
	"sort"
	"util/datafile"
)

const (
	asteroidHit  = '#'
	asteroidMiss = '.'
)

type point struct {
	x, y int
}

type vector struct {
	x, y float64
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (p point) ToVector() vector {
	return vector{float64(p.x), float64(p.y)}
}

func (p point) Subtract(po point) point {
	return point{p.x - po.x, p.y - po.y}
}

func (p point) Manhattan() int {
	return absInt(p.x) + absInt(p.y)
}

func (v vector) Magnitude() float64 {
	return math.Sqrt(math.Pow(v.x, 2) + math.Pow(v.y, 2))
}

func (v vector) Normalize() vector {
	magnitude := v.Magnitude()
	return vector{v.x / magnitude, v.y / magnitude}
}

type asteroidMap struct {
	positions map[point]bool
	size      point
}

func (asteroids asteroidMap) Print() {
	for y := 0; y < asteroids.size.y; y++ {
		for x := 0; x < asteroids.size.x; x++ {
			if asteroids.positions[point{x, y}] {
				fmt.Print(string(asteroidHit))
			} else {
				fmt.Print(string(asteroidMiss))
			}
		}
		fmt.Println()
	}
}

func getAsteroidMapFromFile() asteroidMap {
	file := datafile.Open("advent-2019/day10.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)

	asteroids := asteroidMap{map[point]bool{}, point{0, 0}}
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		asteroids.size.x = len(line)
		for x, c := range line {
			if c == asteroidHit {
				asteroids.positions[point{x, y}] = true
			}
		}
		y++
	}
	asteroids.size.y = y
	return asteroids
}

func part1() {
	asteroids := getAsteroidMapFromFile()
	highestTotal := -1
	highestPos := point{-1, -1}
	for source := range asteroids.positions {
		slopes := map[float64]bool{}
		for dest := range asteroids.positions {
			if source == dest {
				continue
			}

			deltaX := dest.x - source.x
			deltaY := dest.y - source.y
			//dir := point{deltaX, deltaY}.ToVector().Normalize()
			dir := math.Atan2(float64(deltaY), float64(deltaX))
			slopes[dir] = true
		}

		total := len(slopes)
		if highestTotal == -1 || total > highestTotal {
			highestTotal = total
			highestPos = source
		}
	}
	fmt.Println(highestTotal, highestPos)
}

func getLaserShiftedAngle(rad float64) float64 {
	baseAngleInDeg := rad * 180. / math.Pi
	if baseAngleInDeg < 0 {
		baseAngleInDeg += 360
	}
	return math.Mod(baseAngleInDeg + 90, 360)
}

func part2() {
	asteroids := getAsteroidMapFromFile()

	highestTotal := -1
	stationPos := point{-1, -1}
	var highestPosSlopes map[float64][]point
	for station := range asteroids.positions {
		slopes := map[float64][]point{}
		for asteroid := range asteroids.positions {
			if station == asteroid {
				continue
			}

			deltaX := asteroid.x - station.x
			deltaY := asteroid.y - station.y
			dir := math.Atan2(float64(deltaY), float64(deltaX))
			if slopes[dir] == nil {
				slopes[dir] = []point{}
			}
			slopes[dir] = append(slopes[dir], asteroid)
		}

		total := len(slopes)
		if highestTotal == -1 || total > highestTotal {
			highestTotal = total
			stationPos = station
			highestPosSlopes = slopes
		}
	}
	fmt.Println(stationPos)

	var allAngles []float64
	asteroidsByAngle := map[float64][]point{}
	for slopeInRad, points := range highestPosSlopes {
		shiftedAngleInDeg := getLaserShiftedAngle(slopeInRad)
		sort.Slice(points, func(i, j int) bool {
			return stationPos.Subtract(points[i]).Manhattan() < stationPos.Subtract(points[j]).Manhattan()
		})
		asteroidsByAngle[shiftedAngleInDeg] = points
		allAngles = append(allAngles, shiftedAngleInDeg)
		//fmt.Println(shiftedAngleInDeg, asteroidsByAngle[shiftedAngleInDeg])
	}
	sort.Float64s(allAngles)

	vaporized := 0
	vaporizedThisRotation := true
	targetAsteroid := 200
	for vaporizedThisRotation {
		vaporizedThisRotation = false
		for _, angle := range allAngles {
			if asteroidsByAngle[angle] != nil && len(asteroidsByAngle[angle]) > 0 {
				vaporized++
				vaporizedThisRotation = true
				asteroidsAtAngle := asteroidsByAngle[angle]
				asteroidToVaporize := asteroidsAtAngle[0]
				fmt.Println("Vaporizing asteroid", vaporized, "at", asteroidToVaporize, "with angle", angle)
				asteroidsByAngle[angle] = asteroidsAtAngle[1:]

				if vaporized == targetAsteroid {
					fmt.Println("Solution:", asteroidToVaporize.x * 100 + asteroidToVaporize.y)
					return
				}
			}
		}
	}
}

func main() {
	part2()
}
