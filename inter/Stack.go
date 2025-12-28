package inter

type Stack struct {
	s []*ButterVariants
} //定义栈类型Stack

func (s *Stack) size() int64 {
	return int64(len(s.s))
}

func (s *Stack) push(v *ButterVariants) {
	s.s = append(s.s, v)
}

func (s *Stack) pop() {
	s.s = s.s[:len(s.s)-1]
}

func (s *Stack) top() *ButterVariants {
	return s.s[len(s.s)-1]
}

func (s *Stack) bottom() *ButterVariants {
	return s.s[0]
}

func (s *Stack) isEmpty() bool {
	return len(s.s) == 0
}
