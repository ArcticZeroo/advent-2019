package main

import (
	"advent-2019/smath"
	"advent-2019/vector3"
	"fmt"
)

type Moon struct {
	pos vector3.Vector
	vel vector3.Vector
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func applyGravityComponent(posA int, posB int, velA *int, velB *int) {
	if posA == posB {
		return
	}

	if posA > posB {
		*velA--
		*velB++
	} else {
		*velA++
		*velB--
	}
}

func applyGravity(moons []Moon) {
	for i := 0; i < len(moons)-1; i++ {
		moonA := &moons[i]
		for j := i + 1; j < len(moons); j++ {
			moonB := &moons[j]
			applyGravityComponent(moonA.pos.X, moonB.pos.X, &moonA.vel.X, &moonB.vel.X)
			applyGravityComponent(moonA.pos.Y, moonB.pos.Y, &moonA.vel.Y, &moonB.vel.Y)
			applyGravityComponent(moonA.pos.Z, moonB.pos.Z, &moonA.vel.Z, &moonB.vel.Z)
		}
	}
}

func (m *Moon) Move() {
	m.pos = m.pos.Add(m.vel)
}

func (m Moon) PotentialEnergy() int {
	return absInt(m.pos.X) + absInt(m.pos.Y) + absInt(m.pos.Z)
}

func (m Moon) KineticEnergy() int {
	return absInt(m.vel.X) + absInt(m.vel.Y) + absInt(m.vel.Z)
}

func (m Moon) TotalEnergy() int {
	return m.PotentialEnergy() * m.KineticEnergy()
}

func createMoon(pos vector3.Vector) Moon {
	return Moon{pos, vector3.Vector{}}
}

func part1() {
	moons := []Moon{
		createMoon(vector3.Vector{X: -9, Y: -1, Z: -1}),
		createMoon(vector3.Vector{X: 2, Y: 9, Z: 5}),
		createMoon(vector3.Vector{X: 10, Y: 18, Z: -12}),
		createMoon(vector3.Vector{X: -6, Y: 15, Z: -7}),
	}

	for step := 0; step < 1000; step++ {
		applyGravity(moons)
		for i := 0; i < len(moons); i++ {
			moons[i].Move()
			//fmt.Println("Step", step, "moon:", moons[i])
		}
	}

	totalEnergy := 0
	for _, moon := range moons {
		fmt.Println(moon)
		totalEnergy += moon.TotalEnergy()
	}

	fmt.Println("Total energy:", totalEnergy)
}

func part2() {
	moons := []Moon{
		createMoon(vector3.Vector{X: -9, Y: -1,  Z: -1}),
		createMoon(vector3.Vector{X: 2,  Y: 9,   Z: 5}),
		createMoon(vector3.Vector{X: 10, Y: 18, Z: -12}),
		createMoon(vector3.Vector{X: -6, Y: 15,  Z: -7}),
	}

	initialMoons := make([]Moon, len(moons))
	copy(initialMoons, moons)

	cycles := vector3.Vector{-1, -1, -1}
	steps := 0
	for cycles.X == -1 || cycles.Y == -1 || cycles.Z == -1 {
		steps++
		applyGravity(moons)
		matchX, matchY, matchZ := true, true, true
		for i := 0; i < len(moons); i++ {
			moons[i].Move()
			//fmt.Println("Step", step, "moon:", moons[i])
			if moons[i].pos.X != initialMoons[i].pos.X || moons[i].vel.X != 0 {
				matchX = false
			}
			if moons[i].pos.Y != initialMoons[i].pos.Y || moons[i].vel.Y != 0 {
				matchY = false
			}
			if moons[i].pos.Z != initialMoons[i].pos.Z || moons[i].vel.Z != 0 {
				matchZ = false
			}
		}
		if cycles.X == -1 && matchX {
			cycles.X = steps
		}

		if cycles.Y == -1 && matchY {
			cycles.Y = steps
		}

		if cycles.Z == -1 && matchZ {
			cycles.Z = steps
		}
	}

	fmt.Println(steps)
	fmt.Println(smath.LCM(cycles.X, cycles.Y, cycles.Z))
}

func main() {
	part2()
}
