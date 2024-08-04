package config

import (
	"flag"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	CfgType           string         `mapstructure:"cfg_type"`
	PgUrl             string         `mapstructure:"pg_url"`
	MigrationsDir     string         `mapstructure:"migrations_dir"`
	HTTPServerAddress string         `mapstructure:"http_server_address"`
	AuthGPRC          AuthGPRCConfig `mapstructure:"auth_grpc"`
}

type AuthGPRCConfig struct {
	Address   string `mapstructure:"address"`
	SecretKey string `mapstructure:"secret_key"`
	AppId     int32 `mapstructure:"app_id"`
}

var Cfg *Config

func MustLoad() {
	configType := fetchConfigType()
	if configType == "" {
		panic("config type is empty")
	}

	// Reading config from a file
	viper.SetConfigName(configType)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic("Error loading config: " + err.Error())
	}

	// Unmarshalling config to a struct
	config := &Config{}
	err = viper.UnmarshalExact(&config)
	if err != nil {
		panic("Error unmarshalling config to struct: " + err.Error())
	}
	slog.Info("cfg_type: " + config.CfgType)

	Cfg = config
}

// fetchConfigType запрашивает путь к конфигу через командную строку или
// переменную окружения CFG_TYPE.
// Приоритет: flag > env.
// По умолчанию - пустая строка.
func fetchConfigType() string {
	var res string

	flag.StringVar(&res, "config", "", "config type")
	flag.Parse()

	if res == "" {
		// Loading .env file which contains only CFG_TYPE
		err := godotenv.Load()
		if err != nil {
			panic("loading .env file: " + err.Error())
		}
		res = os.Getenv("CFG_TYPE")
	}

	return res
}
