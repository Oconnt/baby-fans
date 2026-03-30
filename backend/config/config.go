package config

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	App    AppConfig    `mapstructure:"app"`
	DB     DBConfig     `mapstructure:"db"`
	WeChat WeChatConfig `mapstructure:"wechat"`
	JWT    JWTConfig    `mapstructure:"jwt"`
	Aliyun AliyunConfig `mapstructure:"aliyun"`
}

type AppConfig struct {
	Version     string `mapstructure:"version"`
	BuildNumber string `mapstructure:"build_number"`
	UpdateURL   string `mapstructure:"update_url"`
	ForceUpdate bool   `mapstructure:"force_update"`
}

type ServerConfig struct {
	Port        string `mapstructure:"port"`
	Mode        string `mapstructure:"mode"`
	Domain      string `mapstructure:"domain"`
	CertDir     string `mapstructure:"cert_dir"`
	Email       string `mapstructure:"email"`
	DNSProvider string `mapstructure:"dns_provider"` // aliyun, cloudflare, etc.
}

type DBConfig struct {
	Type     string `mapstructure:"type"`     // mysql or sqlite
	Host     string `mapstructure:"host"`     // MySQL host
	Port     int    `mapstructure:"port"`     // MySQL port
	Username string `mapstructure:"username"` // MySQL username
	Password string `mapstructure:"password"` // MySQL password
	Name     string `mapstructure:"name"`     // MySQL database name
	DSN      string `mapstructure:"dsn"`      // Custom DSN (overrides other settings)
}

// GetMySQLDSN returns MySQL DSN based on config
func (d *DBConfig) GetMySQLDSN() string {
	if d.DSN != "" {
		return d.DSN
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.Username, d.Password, d.Host, d.Port, d.Name)
}

// GetSQLiteDSN returns SQLite DSN
func (d *DBConfig) GetSQLiteDSN() string {
	if d.DSN != "" {
		return d.DSN
	}
	return "baby-fans.db"
}

type WeChatConfig struct {
	AppID     string `mapstructure:"app_id"`
	AppSecret string `mapstructure:"app_secret"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"` // in hours
}

type AliyunConfig struct {
	AccessKeyID     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
}

var Cfg *Config
var configFilePath string

func init() {
	flag.StringVar(&configFilePath, "config", "", "Path to config file")
}

func LoadConfig() {
	if configFilePath != "" {
		viper.SetConfigFile(configFilePath)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./config")
		viper.AddConfigPath("./etc")
		viper.AddConfigPath("../config")
		viper.AddConfigPath("../etc")
	}

	// Allow overriding via environment variables (e.g., SERVER_PORT)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Default values
	viper.SetDefault("server.port", "18081")
	viper.SetDefault("server.domain", "occont.asia")
	viper.SetDefault("server.cert_dir", "certs")
	viper.SetDefault("server.email", "admin@occont.asia")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.build_number", "100")
	viper.SetDefault("app.update_url", "")
	viper.SetDefault("app.force_update", false)
	viper.SetDefault("jwt.secret", "super_secret_baby_fans_key")
	viper.SetDefault("jwt.expire", 24)
	viper.SetDefault("db.type", "sqlite")
	viper.SetDefault("db.host", "localhost")
	viper.SetDefault("db.port", 3306)
	viper.SetDefault("db.username", "root")
	viper.SetDefault("db.name", "baby_fans")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Error reading config file, using defaults. Error: %v", err)
	}

	Cfg = &Config{}
	if err := viper.Unmarshal(Cfg); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	log.Println("Config loaded successfully")
}

// ReloadConfig reloads config from file without restart
func ReloadConfig() error {
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	newCfg := &Config{}
	if err := viper.Unmarshal(newCfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	Cfg = newCfg
	log.Println("Config reloaded successfully")
	return nil
}
