// Package io contains io abstractions.
package io

import (
	"bufio"
	"io"
	"os"
	"strings"
)

//go:generate mockery --name=WriterFactory

// ClosableWriter ...
type ClosableWriter interface {
	io.Writer
	Close() error
}

// WriterFactory ...
type WriterFactory interface {
	CreateWriter(path string) (ClosableWriter, error)
}

// ClosableStringBuilder ...
type ClosableStringBuilder struct {
	strings.Builder
}

// NewStringWriterFactory returns a new WriterFactory that writes to a string.
func NewStringWriterFactory() WriterFactory {
	return &stringWriterFactory{}
}

// Close ...
func (w *ClosableStringBuilder) Close() error {
	return nil
}

type stringWriterFactory struct{}

// CreateWriter returns a string io.Writer.
func (f *stringWriterFactory) CreateWriter(name string) (ClosableWriter, error) {
	return &ClosableStringBuilder{strings.Builder{}}, nil
	// return &strings.Builder{}, nil
}

type fileWriterFactory struct{}
type closableFileWriter struct {
	bufio.Writer
	file *os.File
}

// NewFileWriterFactory ...
func NewFileWriterFactory() WriterFactory {
	return &fileWriterFactory{}
}

// CreateWriter ...
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

// Close ...
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
