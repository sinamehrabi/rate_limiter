package common

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server struct {
		HTTP struct {
			Host string `yaml:"host" env:"HTTP_HOST"`
			Port string `yaml:"port" env:"HTTP_PORT"`
		} `yaml:"http"`
		TargetHttp struct {
			Host string `yaml:"host" env:"TARGET_HTTP_HOST"`
			Port string `yaml:"port" env:"TARGET_HTTP_PORT"`
		} `yaml:"target_http"`
	} `yaml:"server"`
	Redis struct {
		Host     string `yaml:"host" env:"REDIS_HOST"`
		Port     string `yaml:"port" env:"REDIS_PORT"`
		DB       int    `yaml:"db" env:"REDIS_DB"`
		Username string `yaml:"username" env:"REDIS_USERNAME"`
		Password string `yaml:"password" env:"REDIS_PASSWORD"`
	} `yaml:"redis"`
	Rule struct {
		Path string `yaml:"path"`
		Unit string `yaml:"unit"`
		Rpu  int    `yaml:"rpu"`
	} `yaml:"rule"`
}

func LoadConfig(configPath string) (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err == nil {
		return &cfg, nil
	} else {
		return nil, err
	}
}
