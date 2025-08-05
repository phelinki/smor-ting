package logger

import (
	"context"
	"fmt"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger for easier usage
type Logger struct {
	*zap.Logger
}

// New creates a new logger instance
func New(level, format, output string) (*Logger, error) {
	var config zap.Config

	// Parse log level
	logLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level %s: %w", level, err)
	}

	// Configure based on format
	switch format {
	case "json":
		config = zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(logLevel)
		config.OutputPaths = []string{output}
		config.ErrorOutputPaths = []string{output}
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	case "console":
		config = zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(logLevel)
		config.OutputPaths = []string{output}
		config.ErrorOutputPaths = []string{output}
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	default:
		return nil, fmt.Errorf("unsupported log format: %s", format)
	}

	// Create logger
	zapLogger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	return &Logger{zapLogger}, nil
}

// WithContext adds context fields to the logger
func (l *Logger) WithContext(ctx context.Context) *Logger {
	// Extract request ID from context if available
	if requestID, ok := ctx.Value("request_id").(string); ok {
		return &Logger{l.Logger.With(zap.String("request_id", requestID))}
	}
	return l
}

// WithFields adds multiple fields to the logger
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return &Logger{l.Logger.With(zapFields...)}
}

// Info logs an info level message
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

// Error logs an error level message
func (l *Logger) Error(msg string, err error, fields ...zap.Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	l.Logger.Error(msg, fields...)
}

// Warn logs a warning level message
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, fields...)
}

// Debug logs a debug level message
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, fields...)
}

// Fatal logs a fatal level message and exits
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.Logger.Fatal(msg, fields...)
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

// Default logger instance
var defaultLogger *Logger

// InitDefault initializes the default logger
func InitDefault(level, format, output string) error {
	var err error
	defaultLogger, err = New(level, format, output)
	return err
}

// GetDefault returns the default logger instance
func GetDefault() *Logger {
	if defaultLogger == nil {
		// Fallback to standard logger if not initialized
		log.Println("WARNING: Default logger not initialized, using fallback")
		fallbackLogger, _ := New("info", "console", "stdout")
		return fallbackLogger
	}
	return defaultLogger
}

// Helper functions for easy logging

// Info logs an info message using the default logger
func Info(msg string, fields ...zap.Field) {
	GetDefault().Info(msg, fields...)
}

// Error logs an error message using the default logger
func Error(msg string, err error, fields ...zap.Field) {
	GetDefault().Error(msg, err, fields...)
}

// Warn logs a warning message using the default logger
func Warn(msg string, fields ...zap.Field) {
	GetDefault().Warn(msg, fields...)
}

// Debug logs a debug message using the default logger
func Debug(msg string, fields ...zap.Field) {
	GetDefault().Debug(msg, fields...)
}

// Fatal logs a fatal message and exits using the default logger
func Fatal(msg string, fields ...zap.Field) {
	GetDefault().Fatal(msg, fields...)
}

// Sync flushes the default logger
func Sync() error {
	return GetDefault().Sync()
}
