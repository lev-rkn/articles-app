package logger

import (
	"io"
	"log/slog"
	"os"
)

func Init() {
	// Логирование все файлов в файл logs.txt и в консоль. Установка этого логгера по умолчанию
	outfile, err := os.Create("logs.txt")
	if err != nil {
		slog.Error("unable to create log file", "err", err.Error())
	}
	logger := slog.New(slog.NewJSONHandler(io.MultiWriter(outfile, os.Stdout) , nil))
	slog.SetDefault(logger)
}
