package rrwriter

import (
	"io"
)

type Option func(w *RRWriter) error

type newWriterFunc func(path string) (io.WriteCloser, error)

func WithNewWriter(f newWriterFunc) Option {
	return func(w *RRWriter) error {
		w.newWriter = f
		return nil
	}
}
