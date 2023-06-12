package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
)

func InitLogger() {
	// Configure logger options
	config := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:    "level",
			TimeKey:     "time",
			MessageKey:  "message",
			EncodeLevel: zapcore.CapitalColorLevelEncoder,
			EncodeTime:  zapcore.ISO8601TimeEncoder,
		},
	}

	var err error

	// Create logger instance
	logger, err = config.Build()
	if err != nil {
		panic(err)
	}

	// Create a SugaredLogger for optional sugar syntax
	sugar = logger.Sugar()
}

// Info logs an info-level message
func Info(args ...any) {
	sugar.Info(args...)
}

// Warn logs a warning-level message
func Warn(args ...any) {
	sugar.Warn(args...)
}

// Error logs an error-level message
func Error(args ...any) {
	sugar.Error(args...)
}

// Infof logs a formatted info-level message
func Infof(template string, args ...any) {
	sugar.Infof(template, args...)
}

// Warnf logs a formatted warning-level message
func Warnf(template string, args ...any) {
	sugar.Warnf(template, args...)
}

// Errorf logs a formatted error-level message
func Errorf(template string, args ...any) {
	sugar.Errorf(template, args...)
}

// Sync flushes any buffered logs
func Sync() error {
	return logger.Sync()
}
