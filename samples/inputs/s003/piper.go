package s003

// @pipe
type PipeWithTwoFunctions struct {
	a func (int) float32
	b func (float32) (int, int)
}

// @pipe
type PipeWithThreeFunctions struct {
	a func (int) float32
	b func (float32) (int, int)
	c func (int, int) float64
}

