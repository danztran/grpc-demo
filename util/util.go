package util

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

// StopSignal return a canceling signals channel (like INT, TERM)
func StopSignal() chan os.Signal {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	return quit
}

func Getenv(key string, def string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}

	return def
}

func NewLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	return logger.Sugar()
}
