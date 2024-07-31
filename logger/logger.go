package logger

import (
	"io"
	"log/slog"
	"os"
)

func MustLoad(cfgType string) {
	// writing logs to a logs.txt file
	outfile, err := os.Create("logs.txt")
	if err != nil {
		panic("creating file for logs: " + err.Error())
	}

	var logger *slog.Logger

	switch cfgType {
	case "local":
		logger = slog.New(slog.NewTextHandler(io.MultiWriter(outfile, os.Stdout),
			&slog.HandlerOptions{Level: slog.LevelDebug}))
	case "dev":
		logger = slog.New(slog.NewJSONHandler(io.MultiWriter(outfile, os.Stdout),
			&slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		panic("logger, unknown cfg_type: " + cfgType)
	}

	slog.SetDefault(logger)
}
