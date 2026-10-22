package logger

import (
	"log/slog"
	"os"
)

// Init بدون نیاز به ورودی، محیط را خودش تشخیص می‌دهد
func Init() {
	env := os.Getenv("APP_ENV")
	
	var handler slog.Handler
	// اگر در محیط پروداکشن بودیم یا محیط ست نشده بود، لاگ‌ها JSON و صنعتی باشند
	if env == "production" || env == "" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		// در محیط توسعه (dev)، لاگ‌ها خوانا و متنی باشند
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}
	
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func Info(msg string, args ...any)  { slog.Info(msg, args...) }
func Warn(msg string, args ...any)  { slog.Warn(msg, args...) }
func Error(msg string, args ...any) { slog.Error(msg, args...) }
func Debug(msg string, args ...any) { slog.Debug(msg, args...) }
