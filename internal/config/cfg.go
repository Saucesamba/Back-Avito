package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

const CONFIG_FILE = "../../config.yaml"

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"database"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type AppConfig struct {
	Database DBConfig     `yaml:"database"`
	Server   ServerConfig `yaml:"server"`
}

func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{}
	file, err := os.ReadFile(CONFIG_FILE)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, cfg)
	return cfg, err
}
