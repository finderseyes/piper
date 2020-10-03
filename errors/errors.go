// Package errors contains error-handling utilities
package errors

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
)

var errSkipped = fmt.Errorf("skipping pipeline")

// pipeError contains errors
type pipeError struct {
	error
	stage string
}

// NewError returns a new pipe error.
func NewError(stage string, err error) error {
	return &pipeError{
		error: err,
		stage: stage,
	}
}

// CannotSkip returns an error signifying the stage cannot be skipped.
func CannotSkip(stage string, err error) error {
	return &pipeError{
		error: errors.Wrap(err, "cannot skip"),
		stage: stage,
	}
}

// Skip returns an error that signifies piper that the remaining stages of the pipe should be skipped.
func Skip() error {
	return errSkipped
}

// IsSkipped checks if an error is a skipping signifier.
func IsSkipped(err error) bool {
	return err == errSkipped
}

// IsPipeError checks if an error is piper error.
func IsPipeError(err error) bool {
	_, ok := err.(*pipeError)
	return ok
}

// Stage returns the stage of a piper error. If the error is not from Piper, the function returns empty string.
func Stage(err error) string {
	if pe, ok := err.(*pipeError); ok {
		return pe.Stage()
	}
	return ""
}

// Cause returns the cause of a piper error.
func (w *pipeError) Cause() error { return errors.Cause(w.error) }

// Unwraps returns the cause of a piper error.
func (w *pipeError) Unwrap() error { return w.error }

// Stage returns the stage of a piper error.
func (w *pipeError) Stage() string { return w.stage }

// Format formats a piper error.
func (w *pipeError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = io.WriteString(s, "pipe error in stage: "+w.Stage())
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
