package s005

import (
	errors2 "github.com/finderseyes/piper/errors"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPipeError(t *testing.T) {
	const floatValue = float32(10.0)

	p := &PipeWithError{a: func(i int) (float32, error) {
		return floatValue, errors.New("some error")
	}}

	output, err := p.Run(1)
	assert.Equal(t, output, float32(0.0))
	assert.Error(t, err)
	assert.True(t, errors2.IsPipeError(err))
	assert.Equal(t, errors2.Stage(err), "a")
}
