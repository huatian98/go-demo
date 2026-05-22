package config

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Metrics  MetricsConfig
}

type AppConfig struct {
	Name string
	Env  string
	Port int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	Timezone string
}

type JWTConfig struct {
	Secret      string
	ExpireHours int `mapstructure:"expire_hours"`
}

type MetricsConfig struct {
	IntervalMinutes int `mapstructure:"interval_minutes"`
}

var C *Config

func Load() *Config {
	v := viper.New()
	v.SetConfigType("yaml")

	candidates := []string{
		"./configs/config.yaml",
		"/opt/go-demo/configs/config.yaml",
		"./config.yaml",
	}
	if env := os.Getenv("GO_DEMO_CONFIG"); env != "" {
		candidates = append([]string{env}, candidates...)
	}

	var loaded string
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			v.SetConfigFile(p)
			if err := v.ReadInConfig(); err == nil {
				loaded = p
				break
			}
		}
	}
	if loaded == "" {
		log.Fatalf("config file not found in: %v", candidates)
	}

	v.SetEnvPrefix("GO_DEMO")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		log.Fatalf("config unmarshal failed: %v", err)
	}
	log.Printf("config loaded from %s", loaded)
	C = &c
	return &c
}
