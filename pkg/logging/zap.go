package logging

import (
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newZapLogger(logLevel string, opts ...zap.Option) (*zap.Logger, error) {
	l, err := getLogLevel(logLevel)
	if err != nil {
		return nil, err
	}

	level := zap.NewAtomicLevel()
	level.SetLevel(l)

	config := zap.Config{
		Level:             level,
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     getDefaultEncoderConfig(),
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		InitialFields:     map[string]interface{}{},
	}

	logger, err := config.Build(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to build a Zap logger: %w", err)
	}

	return logger, nil
}


func getDefaultEncoderConfig() zapcore.EncoderConfig {
	// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#logseverity
	return zapcore.EncoderConfig{
		MessageKey:       "msg",
		LevelKey:         "severity",
		TimeKey:          "timestamp",
		NameKey:          "name",
		CallerKey:        "caller",
		FunctionKey:      "",
		StacktraceKey:    "",
		LineEnding:       "",
		EncodeLevel:      zapcore.CapitalLevelEncoder,
		EncodeTime:       zapcore.ISO8601TimeEncoder,
		EncodeDuration:   zapcore.StringDurationEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		EncodeName:       nil,
		ConsoleSeparator: "",
	}
}

var errUnsupportedLogLevel = errors.New("unknown log level")

func getLogLevel(logLevel string) (zapcore.Level, error) {
	var l zapcore.Level

	switch strings.ToLower(logLevel) {
	case "debug":
		l = zapcore.DebugLevel
	case "info":
		l = zapcore.InfoLevel
	case "":
		l = zapcore.InfoLevel
	case "warn":
		l = zapcore.WarnLevel
	case "error":
		l = zapcore.ErrorLevel
	default:
		return l, fmt.Errorf("%s is not supported: %w", logLevel, errUnsupportedLogLevel)
	}

	return l, nil
}
