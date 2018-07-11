package stack

type stack struct {
	data   []interface{}
	length int
}

func (s *stack) Add(item interface{}) {
	if s.length == cap(s.data) {
		newData := make([]interface{}, cap(s.data)*2)
		for i := 0; i < s.length; i++ {
			newData[i] = s.data[i]
		}
		s.data = newData
	}
	s.data[s.length] = item
	s.length++
}

func (s *stack) Remove() interface{} {
	if s.length == 0 {
		panic("underflow")
	}
	s.length--
	item := s.data[s.length]
	s.data[s.length] = nil
	return item
}

func (s *stack) Peek(n int) interface{} {
	if s.length == 0 || n < 0 || n >= s.length {
		panic("range error")
	}
	return s.data[s.length-1-n]
}

func (s *stack) Clear() {
	for i := 0; i < s.length; i++ {
		s.data[i] = nil
	}
	s.length = 0
}

func (s *stack) Len() int {
	return s.length
}

func (s *stack) Cap() int {
	return cap(s.data)
}

func (s *stack) Empty() bool {
	return s.length == 0
}

func (s *stack) AtCapactiy() bool {
	return s.length == cap(s.data)
}
