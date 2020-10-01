// Package errors contains error-handling utilities
package errors

import (
	"fmt"
	"io"
)

// pipeError contains errors
type pipeError struct {
	error
	stage string
}

// New returns a new pipe error.
func New(stage string, err error) error {
	return &pipeError{
		error: err,
		stage: stage,
	}
}

func (w *pipeError) Cause() error { return w.error }

func (w *pipeError) Unwrap() error { return w.error }

func (w *pipeError) Stage() string { return w.stage }

func (w *pipeError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = io.WriteString(s, "pipe error in stage: " + w.Stage())
			fmt.Fprintf(s, "%+v", w.Cause())
			return
		}
		fallthrough
	case 's':
		_, _ = fmt.Fprintf(s, "pipe error in stage: %s \n%s", w.Stage(), w.Cause())
	case 'q':
		_, _ = fmt.Fprintf(s, "pipe error in stage: %s \n%q", w.Stage(), w.Cause())
	}
}


