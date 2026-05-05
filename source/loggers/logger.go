package loggers

import (
	"log"
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
	file   *os.File
}

func InitializeGlobalLogger(logFilePath string) Logger {
	file, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		file = os.Stdout
		log.Println(err)
	}
	handler := slog.NewJSONHandler(file, &slog.HandlerOptions{AddSource: true})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return Logger{logger: logger, file: file}
}

func (l *Logger) Close() error {
	return l.file.Close()
}
