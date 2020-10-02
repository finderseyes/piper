package s000

// @pipe
type SimplePipe struct {
	a func(int) int
}

func NewSimplePipe(a func(int) int) *SimplePipe {
	return &SimplePipe{a: a}
}
