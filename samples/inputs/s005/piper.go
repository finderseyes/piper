package s005

// @pipe
type PipeWithError struct {
	a func(int) (float32, error)
}
