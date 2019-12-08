package rlogging

import (
	"go.uber.org/zap"
	"io"
)

type Logger interface {
	Infof(tmpl string, args ...interface{})
	Debugf(tmpl string, args ...interface{})
	Errorf(tmpl string, args ...interface{})
}

func NewZapLogger() Logger {
	logger, _ := zap.NewProduction()

	defer logger.Sync()

	return &ZapLogger{
		SugaredLogger: logger.Sugar(),
	}
}

func NewWriterLogger(logger Logger) io.Writer {
	return &WriterLogger{
		Logger: logger,
	}
}
