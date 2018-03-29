package logger

import (
	"io"
	"io/ioutil"
)

type Logger interface {
	Verbose() io.Writer
	Info() io.Writer
	Success() io.Writer
	Warning() io.Writer
}

type logger struct {
	verbose io.Writer
	info    io.Writer
	success io.Writer
	warning io.Writer
}

func NewLoggerWriters(v, i, s, w io.Writer) Logger {
	return &logger{
		verbose: v,
		info:    i,
		success: s,
		warning: w,
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

func (l *logger) Warning() io.Writer {
	return discardOrNil(l.warning)
}

func discardOrNil(w io.Writer) io.Writer {
	if w == nil {
		return ioutil.Discard
	}
	return w
}
