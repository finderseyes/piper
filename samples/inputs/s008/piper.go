package s008

// @pipe
type Pipe struct {
	a func(int) (float32, error)
	b func(float32) float32
}

// @pipe
type PipeTwo struct {
	a func(int) (float32, error)
	b func(float32) (float32, error)
}
