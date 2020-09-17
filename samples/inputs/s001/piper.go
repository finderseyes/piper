package s001

// @pipe
type Pipe struct {
	a func(int) int
}

// Not a pipe, should not be generated.
type NotAPipe struct {
	a func(int) int
}
