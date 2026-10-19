package logger
import ("log/slog";"os")
func Init() { slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))) }
func Info(msg string, args ...any) { slog.Info(msg, args...) }
func Error(msg string, args ...any) { slog.Error(msg, args...) }
func Warn(msg string, args ...any) { slog.Warn(msg, args...) }
