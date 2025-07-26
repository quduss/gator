package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}
