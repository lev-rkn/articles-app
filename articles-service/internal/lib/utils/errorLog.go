package utils

import (
	"log/slog"
	"runtime"
)

func ErrorLog(message string, err error) {
	// выводим в дополнение к ошибке текущий файл и линию в коде
	_, file, line, _ := runtime.Caller(1)
	slog.Error(
		message, "err", err.Error(),
		"file", file,
		"line", line)
}