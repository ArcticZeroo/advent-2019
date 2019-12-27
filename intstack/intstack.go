package intstack

type Stack struct {
	data []int
}

func (s Stack) Size() int {
	return len(s.data)
}

func (s Stack) IsEmpty() bool {
	return s.Size() == 0
}

func (s *Stack) Clear() {
	s.data = nil
}

func (s *Stack) Push(x int) {
	s.data = append(s.data, x)
}

func (s *Stack) Pop() int {
	item := s.data[len(s.data) - 1]
	s.data = s.data[:len(s.data) - 1]
	return item
}