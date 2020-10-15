package s000

// @pipe
type P0 struct {
	a func(int) int
}

// @pipe
type P1 struct {
	a func(int)
}

// @pipe
type P2 struct {
	a func(int)
	b func()
}

// @pipe
type P3 struct {
	a func(int) *float32
	b func(*float32) (*float32, error)
}