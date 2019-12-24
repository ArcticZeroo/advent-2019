package point

// Point represents a point on a fixed int grid
type Point struct {
	X, Y int
}

// Subtract returns a point which is the result of the caller point minus the argument
func (p Point) Subtract(po Point) Point {
	return Point{p.X - po.X, p.Y - po.Y}
}

// Add returns a point which is the result of the caller point plus the argument
func (p Point) Add(po Point) Point {
	return Point{p.X + po.X, p.Y + po.Y}
}