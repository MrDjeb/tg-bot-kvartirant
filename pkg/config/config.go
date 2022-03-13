package config

import (
	"os"

	"github.com/spf13/viper"
)

type Text struct {
	Buttons
	Response
}

type Buttons struct {
	Tenant
	Admin
}

type Response struct {
	Start       string `mapstructure:"start"`
	Unknown_cmd string `mapstructure:"unknown_cmd"`
	Unknown_ms  string `mapstructure:"unknown_ms"`
}

type Tenant struct {
	Water1   string `mapstructure:"water1"`
	Receipt1 string `mapstructure:"receipt1"`
	Report1  string `mapstructure:"report1"`
	Water
	Receipt
}

type Admin struct {
	Rooms1 string `mapstructure:"rooms1"`
}

type Water struct {
	Hot_w2  string `mapstructure:"hot_w2"`
	Cold_w2 string `mapstructure:"cold_w2"`
}

type Receipt struct {
	Add_month2   string `mapstructure:"add_month2"`
	Add_amount2  string `mapstructure:"add_amount2"`
	Add_receipt2 string `mapstructure:"add_receipt2"`
}

type Config struct {
	TgToken string
	Text
}

func Init() (*Config, error) {
	if err := setUpViper(); err != nil {
		return nil, err
	}

	var cfg Config

	if err := fromEnv(&cfg); err != nil {
		return nil, err
	}

	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func unmarshal(cfg *Config) error {
	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("text.buttons.tenant", &cfg.Text.Buttons.Tenant); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("text.buttons.tenant.water", &cfg.Text.Buttons.Tenant.Water); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("text.buttons.tenant.receipt", &cfg.Text.Buttons.Tenant.Receipt); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("text.buttons.admin", &cfg.Text.Buttons.Admin); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("text.response", &cfg.Text.Response); err != nil {
		return err
	}

	return nil
}

func fromEnv(cfg *Config) error {
	os.Setenv("TOKEN", "5150854501:AAHM8auF6KgpeHIbw2BHSVMJ5CRPshzYU5s")

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
