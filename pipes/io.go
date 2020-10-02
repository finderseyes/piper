package pipes

import (
	"bufio"
	"io"
	"os"
	"strings"
)

//go:generate mockery --name=WriterFactory

type ClosableWriter interface {
	io.Writer
	Close() error
}

// WriterFactory ...
type WriterFactory interface {
	CreateWriter(path string) (ClosableWriter, error)
}

type closableStringWriter struct {
	strings.Builder
}

// NewStringWriterFactory returns a new WriterFactory that writes to a string.
func NewStringWriterFactory() WriterFactory {
	return &stringWriterFactory{}
}

func (w *closableStringWriter) Close() error {
	return nil
}

type stringWriterFactory struct{}

// CreateWriter returns a string io.Writer.
func (f *stringWriterFactory) CreateWriter(name string) (ClosableWriter, error) {
	return &closableStringWriter{Builder: *&strings.Builder{}}, nil
	// return &strings.Builder{}, nil
}

type fileWriterFactory struct{}
type closableFileWriter struct {
	bufio.Writer
	file *os.File
}

func NewFileWriterFactory() WriterFactory {
	return &fileWriterFactory{}
}

func (f *fileWriterFactory) CreateWriter(path string) (ClosableWriter, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	fileWriter := bufio.NewWriter(file)

	return &closableFileWriter{
		Writer: *fileWriter,
		file:   file,
	}, err
}

func (w *closableFileWriter) Close() error {
	//if err := w.File.Sync(); err != nil {
	//	return err
	//}
	//
	//if err := w.File.Close(); err != nil {
	//	return err
	//}

	if err := w.Writer.Flush(); err != nil {
		return err
	}

	return w.file.Close()
}
