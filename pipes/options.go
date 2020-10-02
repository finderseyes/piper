package pipes

import "github.com/finderseyes/piper/pipes/io"

// Option ...
type Option func(generator *Generator)

// WithWriterFactory sets the WriterFactory of a Generator.
func WithWriterFactory(factory io.WriterFactory) Option {
	return func(generator *Generator) {
		generator.writerFactory = factory
	}
}
