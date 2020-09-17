package pipes

import (
	"io"
	"strings"
)

//go:generate mockery --name=WriterFactory

// WriterFactory ...
type WriterFactory interface {
	CreateWriter(name string) io.Writer
}

// NewStringWriterFactory returns a new WriterFactory that writes to a string.
func NewStringWriterFactory() WriterFactory {
	return &stringWriterFactory{}
}

type stringWriterFactory struct{}

// CreateWriter returns a string io.Writer.
func (f *stringWriterFactory) CreateWriter(name string) io.Writer {
	return &strings.Builder{}
}
