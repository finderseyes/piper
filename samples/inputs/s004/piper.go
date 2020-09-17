package s004

type Functor interface {
	InvokeByAnotherName(float32) (int, int)
}

// @pipe
type PipeWithFunctor struct {
	a func(int) float32
	b Functor
}
