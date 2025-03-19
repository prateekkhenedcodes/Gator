package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)
const configFileName = ".gatorconfig.json"

type Config struct {

    DBUrl           string `json:"db_url"`
    CurrentUserName string `json:"current_user_name"`
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
	if err != nil{
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

func (c *Config)SetUser(u string)  error{
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	c.CurrentUserName = u
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath, data, 0666)
	if err != nil {
		return err
	}
	return nil 
}