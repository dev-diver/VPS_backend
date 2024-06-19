package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	DBUser string `json:"DBUser"`
	DBPass string `json:"DBPass"`
	DBHost string `json:"DBHost"`
	DBPort string `json:"DBPort"`
	DBName string `json:"DBName"`
}

func LoadConfig(configFile string) (*Config, error) {
	config := &Config{}
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}
	err = json.Unmarshal(file, config)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal config JSON: %w", err)
	}
	return config, nil
}
