package logger

import (
	"log/slog"
	"os"
)

func Init(env string) {
	var handler slog.Handler
	if env == "production" {
		// لاگ‌های صنعتی به صورت JSON برای مانیتورینگ
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		// لاگ‌های توسعه به صورت متنی و با جزئیات بیشتر
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func Info(msg string, args ...any)  { slog.Info(msg, args...) }
func Warn(msg string, args ...any)  { slog.Warn(msg, args...) }
func Error(msg string, args ...any) { slog.Error(msg, args...) }
func Debug(msg string, args ...any) { slog.Debug(msg, args...) }
