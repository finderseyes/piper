package s006

// @pipe
type PipeThree struct {
	a func(int) (float32, error)
	b func(float32) float64
	c func(float64) (int64, error)
}
