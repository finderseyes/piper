package pipes

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/finderseyes/piper/pipes/io"

	"github.com/dave/dst/decorator"

	"github.com/finderseyes/piper/pipes/io/mocks"
	_ "github.com/finderseyes/piper/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewGenerator(t *testing.T) {
	g := NewGenerator("some/path")
	assert.NotNil(t, g)
}

func TestGenerator_Execute_FailedByPath(t *testing.T) {
	var tests = []struct {
		path string
	}{
		{path: "not/exist"},
		{path: "samples/piper-0.go"},
	}

	for _, c := range tests {
		path := c.path
		t.Run(c.path, func(t *testing.T) {
			g := NewGenerator(path)
			err := g.Execute()
			assert.Error(t, err)
		})
	}
}

func TestGenerator_Execute_Succeed(t *testing.T) {
	const count = 8

	for i := 0; i < count; i++ {
		input := fmt.Sprintf("samples/inputs/s%03d", i)
		output := fmt.Sprintf("samples/outputs/s%03d.gen", i)

		t.Run(input, func(t *testing.T) {
			writer := &io.ClosableStringBuilder{}
			mockFactory := &mocks.WriterFactory{}
			mockFactory.On("CreateWriter",
				mock.AnythingOfType("string"),
			).Return(writer, nil)
			g := NewGenerator(input, WithWriterFactory(mockFactory))
			err := g.Execute()

			assert.NoError(t, err)

			buff, _ := ioutil.ReadFile(output)

			resultTree, err := decorator.Parse(writer.String())
			assert.NoError(t, err)
			result := &strings.Builder{}
			_ = decorator.Fprint(result, resultTree)

			expectedTree, err := decorator.Parse(strings.ReplaceAll(string(buff), "\r\n", "\n"))
			assert.NoError(t, err)
			expected := &strings.Builder{}
			_ = decorator.Fprint(expected, expectedTree)

			assert.Equal(t, expected.String(), result.String())
		})
	}
}
