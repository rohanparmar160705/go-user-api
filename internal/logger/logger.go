/*
Package logger initializes the Uber Zap structured logger.
It allows for global access to the logger instance and configures it based on the environment
(e.g., human-readable console logs for development, JSON logs for production).
*/
package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

// InitLogger initializes the Uber Zap logger
func InitLogger(env string) error {
	var err error
	
	if env == "production" {
		Log, err = zap.NewProduction()
	} else {
		Log, err = zap.NewDevelopment()
	}
	
	if err != nil {
		return err
	}
	
	return nil
}

// Sync flushes any buffered log entries
func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}

