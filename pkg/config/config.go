package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	TgToken string
	Text
}
type Text struct {
	Buttons
	Response
}

type Buttons struct {
	Tenant Tenant
	Admin  Admin
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
	Water    Water
	Receipt  Receipt
}

type Water struct {
	Hot_w2  string `mapstructure:"hot_w2"`
	Cold_w2 string `mapstructure:"cold_w2"`
}

type Receipt struct {
	Month2   string `mapstructure:"month2"`
	Amount2  string `mapstructure:"amount2"`
	Receipt2 string `mapstructure:"receipt2"`
}

type Admin struct {
	Rooms1    string `mapstructure:"rooms1"`
	Settings1 string `mapstructure:"settings1"`
	Rooms     Rooms
	Settings  Settings
}

type Rooms struct {
	R2 string `mapstructure:"r2"`
}
type Settings struct {
	Edit2     string `mapstructure:"edit2"`
	Contacts2 string `mapstructure:"contacts2"`
	Reminder2 string `mapstructure:"reminder2"`
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
	if err := viper.UnmarshalKey("text.buttons.admin.rooms", &cfg.Text.Buttons.Admin.Rooms); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("text.buttons.admin.settings", &cfg.Text.Buttons.Admin.Settings); err != nil {
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
