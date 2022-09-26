package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	TgToken string
	Text
}

type Text struct {
	JwtKey     string
	GodModeKey string
	Buttons
	Response Response
	CommonCommand
	CommonMessage
	Constants Constants
}

type Constants struct {
	MAX_SHOW_SCORER  int `mapstructure:"max_show_scorer"`
	MAX_SHOW_PAYMENT int `mapstructure:"max_show_payment"`
}

type Buttons struct {
	Tenant Tenant
	Admin  Admin
}

type Response struct {
	Start          string `mapstructure:"start"`
	Unknown_cmd    string `mapstructure:"unknown_cmd"`
	Unknown_ms     string `mapstructure:"unknown_ms"`
	Cache_ttl      string `mapstructure:"cache_ttl"`
	Water1_first   string `mapstructure:"water1_first"`
	Water1_sec     string `mapstructure:"water1_sec"`
	Water2_inp     string `mapstructure:"water2_inp"`
	Water2_change  string `mapstructure:"water2_change"`
	Water2_saved   string `mapstructure:"water2_saved"`
	Receipt1_first string `mapstructure:"receipt1_first"`
	Receipt1_sec   string `mapstructure:"receipt1_sec"`
	Receipt1_third string `mapstructure:"receipt1_third"`
	Amount2_inp    string `mapstructure:"amount2_inp"`
	Receipt2_saved string `mapstructure:"receipt2_saved"`
	Report1_info   string `mapstructure:"report1_info"`
	Rooms1_list    string `mapstructure:"rooms1_list"`
	Rooms1_nil     string `mapstructure:"rooms1_nil"`
	Settings1      string `mapstructure:"settings1"`
	Room2          string `mapstructure:"room2"`
}

type CommonCommand struct {
	Start   string `mapstructure:"start"`
	Cancel  string `mapstructure:"cancel"`
	Unknown string `mapstructure:"unknown"`
	BackBut string `mapstructure:"back_but"`
	GodMode string `mapstructure:"godmode"`
}

type CommonMessage struct {
	Hi string `mapstructure:"hi"`
}
type Tenant struct {
	Water1   string `mapstructure:"water1"`
	Receipt1 string `mapstructure:"receipt1"`
	Report1  string `mapstructure:"report1"`
	Water    Water
	Receipt  Receipt
}

type Water struct {
	Hot_w2       string `mapstructure:"hot_w2"`
	Cold_w2      string `mapstructure:"cold_w2"`
	Choose_month string `mapstructure:"choose_month"`
	Month_prefix string `mapstructure:"month_prefix"`
}

type Receipt struct {
	Month2       string `mapstructure:"month2"`
	Amount2      string `mapstructure:"amount2"`
	Receipt2     string `mapstructure:"receipt2"`
	Month_prefix string `mapstructure:"month_prefix"`
}

type Admin struct {
	Rooms1    string `mapstructure:"rooms1"`
	Settings1 string `mapstructure:"settings1"`
	Room2     string `mapstructure:"room2"`
	Room      Room
	Settings  Settings
}

type Room struct {
	Payment_prefix string `mapstructure:"payment_prefix"`
	ShowScorer33   string `mapstructure:"show_scorer33"`
	ShowScorerN4   string `mapstructure:"show_scorerN4"`
	ShowScorerB3   string `mapstructure:"show_scorerB3"`
	ShowPayment33  string `mapstructure:"show_payment33"`
	ShowPaymentN4  string `mapstructure:"show_paymentN4"`
	ShowPaymentB3  string `mapstructure:"show_paymentB3"`
	ShowTenants3   string `mapstructure:"show_tenants3"`
}

type Settings struct {
	Edit2         string `mapstructure:"edit2"`
	Contacts2     string `mapstructure:"contacts2"`
	Reminder2     string `mapstructure:"reminder2"`
	ReminderSend3 string `mapstructure:"reminder_send3"`
	ReminderEdit3 string `mapstructure:"reminder_edit3"`
	Edit          Edit
}

type Edit struct {
	AddRoom3    string `mapstructure:"add_room3"`
	RemoveRoom3 string `mapstructure:"remove_room3"`
	Removing4   string `mapstructure:"removing4"`
	Removing    Removing
}

type Removing struct {
	ConfirmRemove5 string `mapstructure:"confirm_remove5"`
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
	if err := viper.UnmarshalKey("text.buttons.admin.room", &cfg.Text.Buttons.Admin.Room); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("text.buttons.admin.settings.edit", &cfg.Text.Buttons.Admin.Settings.Edit); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("text.buttons.admin.settings.edit.removing", &cfg.Text.Buttons.Admin.Settings.Edit.Removing); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("text.buttons.admin.settings", &cfg.Text.Buttons.Admin.Settings); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("text.response", &cfg.Text.Response); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("text.common_command", &cfg.Text.CommonCommand); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("text.common_message", &cfg.Text.CommonMessage); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("text.constants", &cfg.Text.Constants); err != nil {
		return err
	}

	return nil
}

func fromEnv(cfg *Config) error {
	godotenv.Load()

	if err := viper.BindEnv("tg_token"); err != nil {
		return err
	}
	cfg.TgToken = viper.GetString("tg_token")

	if err := viper.BindEnv("access_secret"); err != nil {
		return err
	}
	cfg.Text.JwtKey = viper.GetString("access_secret")

	if err := viper.BindEnv("godmode_secret"); err != nil {
		return err
	}
	cfg.Text.GodModeKey = viper.GetString("godmode_secret")

	return nil
}

func setUpViper() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	return viper.ReadInConfig()
}
