package rrwriter

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync/atomic"
)

// RRWriter write to constant number files in round-robin strategy.
type RRWriter struct {
	fws   []io.WriteCloser
	paths []string

	cnt  int
	line uint32

	closeCallBack func(path string) error
	newWriter     newWriterFunc
}

// NewPathByIndex use for get target filepath by index.
type NewPathByIndex func(idx int) string

// CallOnPath will call on closing each file.
type CallOnPath func(path string) error

func New(
	cnt int,
	getTargetPath NewPathByIndex,
	closeCallBack CallOnPath,
	opts ...Option,
) (*RRWriter, error) {
	if cnt == 0 {
		return nil, errors.New("cnt can not be 0")
	}

	w := &RRWriter{
		cnt:           cnt,
		fws:           make([]io.WriteCloser, cnt),
		paths:         make([]string, cnt),
		closeCallBack: closeCallBack,
		newWriter: func(path string) (io.WriteCloser, error) {
			return os.Create(path)
		},
	}
	for _, o := range opts {
		if err := o(w); err != nil {
			return nil, err
		}
	}

	for i := 0; i < cnt; i++ {
		path := getTargetPath(i)
		fw, err := w.newWriter(path)
		if err != nil {
			return nil, fmt.Errorf("new writer: %w", err)
		}
		w.fws[i] = fw
		w.paths[i] = path
	}
	return w, nil
}

func (w *RRWriter) Close() error {
	var errStr []string
	for _, fw := range w.fws {
		if err := fw.Close(); err != nil {
			errStr = append(errStr, err.Error())
		}
	}
	if len(errStr) != 0 {
		return errors.New(strings.Join(errStr, ","))
	}

	// do callback
	for _, path := range w.paths {
		if err := w.closeCallBack(path); err != nil {
			return fmt.Errorf("close callback: %w", err)
		}
	}

	return nil
}

func (w *RRWriter) Write(p []byte) (int, error) {
	idx := int(atomic.AddUint32(&w.line, 1)) % w.cnt
	return w.fws[idx].Write(p)
}
