package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	DB     DBConfig     `mapstructure:"db"`
	WeChat WeChatConfig `mapstructure:"wechat"`
	JWT    JWTConfig    `mapstructure:"jwt"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DBConfig struct {
	DSN string `mapstructure:"dsn"`
}

type WeChatConfig struct {
	AppID     string `mapstructure:"app_id"`
	AppSecret string `mapstructure:"app_secret"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"` // in hours
}

var Cfg *Config

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config") // Relative to where you run the binary (backend dir)
	viper.AddConfigPath("./etc")    // For etc directory
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../etc")

	// Allow overriding via environment variables (e.g., SERVER_PORT)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Default values
	viper.SetDefault("server.port", "18081")
	viper.SetDefault("jwt.secret", "super_secret_baby_fans_key")
	viper.SetDefault("jwt.expire", 24)

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Error reading config file, using defaults. Error: %v", err)
	}

	Cfg = &Config{}
	if err := viper.Unmarshal(Cfg); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	log.Println("Config loaded successfully")
}
