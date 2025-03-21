package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/prateekkhenedcodes/Gator/internal/database"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) Save() error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	homeDir, err := getConfigFilePath()
	if err != nil {
		return err
	}
	return os.WriteFile(homeDir, data, 0644)
}

type State struct {
	Db        *database.Queries
	ConfigPtr *Config
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	filePath := filepath.Join(homeDir, configFileName)
	return filePath, nil
}
func Read() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	rawData, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}
	config := Config{}
	err = json.Unmarshal(rawData, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}


