package logger

import (
	"io"
	"io/ioutil"
)

type Logger interface {
	Verbose() io.Writer
	Info() io.Writer
	Success() io.Writer
	Error() io.Writer
}

type logger struct {
	verbose io.Writer
	info    io.Writer
	success io.Writer
	err     io.Writer
}

func NewLoggerWriters(v, i, s, e io.Writer) Logger {
	return &logger{
		verbose: v,
		info:    i,
		success: s,
		err:     e,
	}
}

func (l *logger) Verbose() io.Writer {
	return discardOrNil(l.verbose)
}

func (l *logger) Info() io.Writer {
	return discardOrNil(l.info)
}

func (l *logger) Success() io.Writer {
	return discardOrNil(l.success)
}

func (l *logger) Error() io.Writer {
	return discardOrNil(l.err)
}

func discardOrNil(w io.Writer) io.Writer {
	if w == nil {
		return ioutil.Discard
	}
	return w
}
