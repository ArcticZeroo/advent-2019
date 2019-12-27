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

var Directions = []Direction{Up, Right, Down, Left,}

func wrap(v, l int) int {
	for v < 0 {
		v += l
	}

	return v % l
}

func (dir Direction) IsOpposite(otherDir Direction) bool {
	return dir.Opposite() == otherDir
}

func (dir Direction) Opposite() Direction {
	switch dir {
	case Up:
		return Down
	case Down:
		return Up
	case Left:
		return Right
	case Right:
		return Left
	}
	return -1
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
	t.pos = NextPosInDir(t.pos, dir)
}

func (t Turtle) Pos() point.Point {
	return t.pos
}

// NextPosInDir returns a position which is exactly 1 grid point away from the given pos, in the given direction
func NextPosInDir(pos point.Point, dir Direction) point.Point {
	component := &pos.X
	if dir == Up || dir == Down {
		component = &pos.Y
	}
	increment := 1
	if dir == Down || dir == Left {
		increment = -1
	}
	*component += increment
	return pos
}