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

func NewZapLogger(options ...func(zap.Config) zap.Config) (Logger, error) {
	c := zap.NewProductionConfig()
	for _, o := range options {
		c = o(c)
	}

	logger, err := c.Build()
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

func WithOutputPath(outputPath string) func(zap.Config) zap.Config {
	return func(c zap.Config) zap.Config {
		if c.OutputPaths == nil {
			c.OutputPaths = []string{}
		}

		c.OutputPaths = append(c.OutputPaths, outputPath)

		return c
	}
}
