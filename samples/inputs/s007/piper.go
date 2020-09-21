package s000

import (
	"github.com/finderseyes/piper/samples/inputs/s007/child"
	"github.com/finderseyes/piper/samples/inputs/s007/childtwo"
)

// @pipe
type SimplePipe struct {
	a child.Foo
	b func(*child.Data) childtwo.Intptrptr
}
