package pipes

// Option ...
type Option func(generator *Generator)

// WithWriterFactory sets the WriterFactory of a Generator.
func WithWriterFactory(factory WriterFactory) Option {
	return func(generator *Generator) {
		generator.writerFactory = factory
	}
}
