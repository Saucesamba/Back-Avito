package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

const CONFIG_FILE = "config.yaml"

type DBConfig struct {
	Host     string `yaml:"dbhost"`
	Port     string `yaml:"dbport"`
	User     string `yaml:"dbuser"`
	Password string `yaml:"dbpassword"`
	Name     string `yaml:"dbdatabase"`
}

type ServerConfig struct {
	Host string `yaml:"shost"`
	Port string `yaml:"sport"`
}

type AppConfig struct {
	Database *DBConfig    `yaml:"database"`
	Server   ServerConfig `yaml:"server"`
	JWT      JWTConfig    `yaml:"jwt"`
}

type JWTConfig struct {
	TokenExpiryHours int    `yaml:"tokenExpiryHours"`
	Secret           string `yaml:"secret"`
}

func (a *AppConfig) LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{}
	file, err := os.ReadFile(CONFIG_FILE)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, cfg)
	return cfg, err

}
