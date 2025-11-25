package loggers

import (
	"io"
	"log"
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
	file   *os.File
}

func InitializeGlobalLogger() Logger {
	file, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.Panic(err)
	}
	handler := slog.NewJSONHandler(io.MultiWriter(os.Stderr, file), &slog.HandlerOptions{AddSource: true})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return Logger{logger: logger, file: file}
}

func (l *Logger) Close() error {
	return l.file.Close()
}
