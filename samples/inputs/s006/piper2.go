package s006

// @pipe
type PipeTwo struct {
	a func(int) (float32, error)
	b func(float32) float64
}
