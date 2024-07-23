package config

import (
	"log/slog"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func New(logger *slog.Logger) *koanf.Koanf {
	var k = koanf.New(".")
	if err := k.Load(file.Provider("config/conf.yaml"), yaml.Parser()); err != nil {
		logger.Error("loading config",
			"err", err.Error())
	}

	return k
}
