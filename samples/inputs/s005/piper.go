package s005

// @pipe
type PipeWithError struct {
	a func(int) (float32, error)
}

// @pipe
type P1 struct {
	a func(int) error
}

// @pipe
type P2 struct {
	a func(int) (float32, error)
	b func(float32)
}
