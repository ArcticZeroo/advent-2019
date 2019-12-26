package turtle

import "advent-2019/point"

type Direction int

const (
	Up             = 0
	Right          = Up + 1
	Down           = Right + 1
	Left           = Down + 1
	directionCount = Left + 1
)

func wrap(v, l int) int {
	for v < 0 {
		v += l
	}

	return v % l
}

type Turtle struct {
	pos point.Point
	dir Direction
}

func (t *Turtle) MoveForward() {
	t.MoveInDir(t.dir)
}

func (t *Turtle) TurnRight() {
	t.dir = Direction(wrap(int(t.dir + 1), directionCount))
}

func (t *Turtle) TurnLeft() {
	t.dir = Direction(wrap(int(t.dir - 1), directionCount))
}

func (t *Turtle) MoveInDir(dir Direction) {
	component := &t.pos.X
	if dir == Up || dir == Down {
		component = &t.pos.Y
	}
	increment := 1
	if dir == Down || dir == Left {
		increment = -1
	}
	*component += increment
}

func (t Turtle) Pos() point.Point {
	return t.pos
}