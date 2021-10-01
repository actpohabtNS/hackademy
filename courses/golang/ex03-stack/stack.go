package stack

type Stack []int

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack) Push(val int) {
	*s = append(*s, val)
}

func (s *Stack) Pop() int {
	if s.IsEmpty() {
		return 0
	}

	idx := len(*s) - 1
	elem := (*s)[idx]
	*s = (*s)[:idx]
	return elem
}

func New() *Stack {
	return &Stack{}
}
