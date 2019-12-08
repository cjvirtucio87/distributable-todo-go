package rlogging

import "io"

type WriterLogger struct {
	Logger
	io.Writer
}

func (w *WriterLogger) Write(p []byte) (n int, err error) {
	w.Infof(string(p))

	return len(p), nil
}
