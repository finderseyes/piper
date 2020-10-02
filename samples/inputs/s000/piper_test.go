package s000

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestFailed(t *testing.T) {
	p := NewSimplePipe(func(i int) int {
		return i * 2
	})

	assert.Equal(t, p.Run(1), 2)
}
