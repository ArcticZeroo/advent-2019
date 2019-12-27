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

type Queue struct {
	data []Point
}

func (queue *Queue) Push(point Point) {
	queue.data = append(queue.data, point)
}

func (queue *Queue) Pop() Point {
	point := queue.data[0]
	queue.data = queue.data[1:]
	return point
}

func (queue Queue) Size() int {
	return len(queue.data)
}

func (queue Queue) IsEmpty() bool {
	return queue.Size() == 0
}

func (queue *Queue) Clear() {
	queue.data = nil
}

func (queue Queue) Copy() Queue {
	data := make([]Point, len(queue.data))
	copy(data, queue.data)
	return Queue{data: data}
}
