package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func New(logger *slog.Logger) *koanf.Koanf {
	err := godotenv.Load()
	if err != nil {
		logger.Info("loading .env file", "err", err.Error())
	}
	logger.Info("env_type: " + os.Getenv("ENV_TYPE"))

	var k = koanf.New(".")
	configYamlPath := "config/" + os.Getenv("ENV_TYPE") + ".yaml"
	if err := k.Load(file.Provider(configYamlPath), yaml.Parser()); err != nil {
		logger.Error("loading config",
			"err", err.Error())
	}

	return k
}
