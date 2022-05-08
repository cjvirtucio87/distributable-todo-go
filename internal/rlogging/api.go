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

func NewZapLogger() (Logger, error) {
	cfg := zap.NewProductionConfig()
	// TODO: this needs to be configurable
	cfg.OutputPaths = []string{
		"/tmp/test.log",
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &ZapLogger{
		SugaredLogger: logger.Sugar(),
	}, nil
}

func NewWriterLogger(logger Logger) io.Writer {
	return &WriterLogger{
		Logger: logger,
	}
}
