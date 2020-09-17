package pipes

import (
	"io"
	"strings"
)

//go:generate mockery --name=WriterFactory
type WriterFactory interface {
	CreateWriter(name string) io.Writer
}

func NewStringWriterFactory() WriterFactory {
	return &stringWriterFactory{}
}

type stringWriterFactory struct {}

func (f *stringWriterFactory) CreateWriter(name string) io.Writer {
	return &strings.Builder{}
}
