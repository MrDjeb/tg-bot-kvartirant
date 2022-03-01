package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	TgToken string
}

func Init() (*Config, error) {
	if err := setUpViper(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := getEnv(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func getEnv(cfg *Config) error {
	if err := viper.BindEnv("token"); err != nil {
		return err
	}
	cfg.TgToken = viper.GetString("token")
	return nil
}

func setUpViper() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	return viper.ReadInConfig()
}
