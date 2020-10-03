package s008

import (
	"testing"

	perrors "github.com/finderseyes/piper/errors"
	"github.com/stretchr/testify/assert"
)

func TestPipeError(t *testing.T) {
	p := &Pipe{
		a: func(i int) (float32, error) {
			if i > 5 {
				return float32(i * 10), perrors.Skip()
			}
			return float32(i), nil
		},
		b: func(f float32) float32 {
			return f * 2
		},
	}

	{
		output, err := p.Run(10)
		assert.Equal(t, float32(100), output)
		assert.Error(t, err)
		assert.True(t, perrors.IsSkipped(err))
	}

	{
		output, err := p.Run(2)
		assert.Equal(t, float32(4), output)
		assert.NoError(t, err)
	}
}

func TestPipeTwoError(t *testing.T) {
	p := &PipeTwo{
		a: func(i int) (float32, error) {
			if i > 5 {
				return float32(i * 10), perrors.Skip()
			}
			return float32(i), nil
		},
		b: func(f float32) (float32, error) {
			return f * 2, nil
		},
	}

	{
		output, err := p.Run(10)
		assert.Equal(t, float32(100), output)
		assert.Error(t, err)
		assert.True(t, perrors.IsSkipped(err))
	}

	{
		output, err := p.Run(2)
		assert.Equal(t, float32(4), output)
		assert.NoError(t, err)
	}
}
