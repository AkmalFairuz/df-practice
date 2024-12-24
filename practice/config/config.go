package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Server struct {
		Name          string `yaml:"name"`
		ListenAddress string `yaml:"listen_address"`
		MaxPlayers    int    `yaml:"max_players"`
	} `yaml:"server"`

	Database struct {
		Hostname string `yaml:"hostname"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`

	ViewDistance int `yaml:"view_distance"`
}

func (c Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", c.Database.Username, c.Database.Password, c.Database.Hostname, c.Database.Port, c.Database.Name)
}

var globalConfig *Config

func Get() *Config {
	return globalConfig
}

func init() {
	configPath := "config.yml"
	if os.Getenv("CONFIG_PATH") != "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to read config file: %w", err))
	}

	globalConfig = &Config{}
	if err := yaml.Unmarshal(configBytes, globalConfig); err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %w", err))
	}
}
