package dbschema

import "io"

type Logger interface {
	Verbose() io.Writer
	Info() io.Writer
	Error() io.Writer
}

func NewLoggerWriters(verbose, info, err io.Writer) Logger {
	return &logger{
		verbose: verbose,
		info:    info,
		err:     err,
	}
}

type logger struct {
	verbose io.Writer
	info    io.Writer
	err     io.Writer
}

func (l *logger) Verbose() io.Writer {
	return l.verbose
}

func (l *logger) Info() io.Writer {
	return l.info
}

func (l *logger) Error() io.Writer {
	return l.err
}
