package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	LogConfigDefault = zap.Config{
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling: &zap.SamplingConfig{
			Initial:    3,
			Thereafter: 10,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "level",
			TimeKey:        "time",
			NameKey:        "name",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		ErrorOutputPaths: []string{"stderr"},
	}
)

func New(output string) (logger *zap.SugaredLogger, sync func(), err error) {
	c := LogConfigDefault

	level := zap.NewAtomicLevel()
	level.SetLevel(zapcore.InfoLevel)
	c.Level = level

	c.OutputPaths = []string{output}

	l, err := c.Build()
	if err != nil {
		return nil, nil, err
	}

	sync = func() {
		if err := l.Sync(); err != nil && output != "stdout" {
			_, _ = fmt.Fprintln(os.Stderr, "[error]", err)
		}
	}

	logger = l.Sugar()

	return
}
