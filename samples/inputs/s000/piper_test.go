package s000

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestP0_Success(t *testing.T) {
	p := &P0{func(i int) int {
		return i * 2
	}}

	assert.Equal(t, p.Run(1), 2)
}

func TestP1_Success(t *testing.T) {
	received := 0

	p := &P1{a: func(i int) {
		received = i
	}}

	p.Run(1)
	assert.Equal(t, 1, received)
}
