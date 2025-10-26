package utils

import (
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

// InitLogger configures logrus for production use
func InitLogger() {
	Logger = logrus.New()

	// Set formatter based on environment
	if os.Getenv("ENV") == "development" {
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05Z07:00",
			PrettyPrint:     false,
		})
	}

	Logger.SetOutput(os.Stdout)

	// Set log level from environment
	levelStr := os.Getenv("LOG_LEVEL")
	if levelStr == "" {
		levelStr = "info"
	}
	level, err := logrus.ParseLevel(strings.ToLower(levelStr))
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)
}

// GetLogger returns the configured logger instance
func GetLogger() *logrus.Logger {
	if Logger == nil {
		InitLogger()
	}
	return Logger
}

// LogWithFields creates a logger entry with structured fields
func LogWithFields(level logrus.Level, message string, fields logrus.Fields) {
	logger := GetLogger()
	entry := logger.WithFields(fields)

	switch level {
	case logrus.DebugLevel:
		entry.Debug(message)
	case logrus.InfoLevel:
		entry.Info(message)
	case logrus.WarnLevel:
		entry.Warn(message)
	case logrus.ErrorLevel:
		entry.Error(message)
	case logrus.FatalLevel:
		entry.Fatal(message)
	case logrus.PanicLevel:
		entry.Panic(message)
	}
}

// LogError logs an error with context
func LogError(err error, message string, fields ...logrus.Fields) {
	logger := GetLogger()
	entry := logger.WithError(err)

	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}

	entry.Error(message)
}

// LogInfo logs an info message with optional fields
func LogInfo(message string, fields ...logrus.Fields) {
	logger := GetLogger()

	if len(fields) > 0 {
		entry := logger.WithFields(fields[0])
		entry.Info(message)
	} else {
		logger.Info(message)
	}
}

// LogDebug logs a debug message with optional fields
func LogDebug(message string, fields ...logrus.Fields) {
	logger := GetLogger()

	if len(fields) > 0 {
		entry := logger.WithFields(fields[0])
		entry.Debug(message)
	} else {
		logger.Debug(message)
	}
}

// LogWarn logs a warning message with optional fields
func LogWarn(message string, fields ...logrus.Fields) {
	logger := GetLogger()

	if len(fields) > 0 {
		entry := logger.WithFields(fields[0])
		entry.Warn(message)
	} else {
		logger.Warn(message)
	}
}

// LogFatal logs a fatal message and exits
func LogFatal(message string, fields ...logrus.Fields) {
	logger := GetLogger()

	if len(fields) > 0 {
		entry := logger.WithFields(fields[0])
		entry.Fatal(message)
	} else {
		logger.Fatal(message)
	}
}

// LogRequest logs HTTP request information
func LogRequest(method, path, userAgent string, statusCode int, duration time.Duration, fields ...logrus.Fields) {
	logger := GetLogger()
	entry := logger.WithFields(logrus.Fields{
		"method":      method,
		"path":        path,
		"user_agent":  userAgent,
		"status_code": statusCode,
		"duration":    duration.String(),
	})

	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}

	entry.Info("HTTP Request")
}

// LogDatabase logs database operations
func LogDatabase(operation, table string, duration time.Duration, err error, fields ...logrus.Fields) {
	logger := GetLogger()
	entry := logger.WithFields(logrus.Fields{
		"operation": operation,
		"table":     table,
		"duration":  duration.String(),
	})

	if err != nil {
		entry = entry.WithError(err)
	}

	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}

	if err != nil {
		entry.Error("Database operation failed")
	} else {
		entry.Info("Database operation completed")
	}
}

// LogTrade logs trade-related operations
func LogTrade(operation, tradeID, symbol string, userID int, fields ...logrus.Fields) {
	logger := GetLogger()
	entry := logger.WithFields(logrus.Fields{
		"operation": operation,
		"trade_id":  tradeID,
		"symbol":    symbol,
		"user_id":   userID,
	})

	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}

	entry.Info("Trade operation")
}

// GenerateID generates a unique ID
func GenerateID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// GetCurrentTime returns the current time
func GetCurrentTime() time.Time {
	return time.Now()
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
