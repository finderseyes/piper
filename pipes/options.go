package pipes

// Option ...
type Option func (generator *Generator)

func WithWriterFactory(factory WriterFactory) Option {
	return func(generator *Generator) {
		generator.writerFactory = factory
	}
}
