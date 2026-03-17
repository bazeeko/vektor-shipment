package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const envConfigFilePath = "CONFIG_FILE_PATH"

type Config struct {
	Server Server `yaml:"server"`
}

type Server struct {
	GRPCPort int `yaml:"grpc_port"`
}

func Load() (*Config, error) {
	file, err := os.Open(os.Getenv(envConfigFilePath))
	if err != nil {
		return nil, fmt.Errorf("os.Open: %w", err)
	}
	defer file.Close()

	var config Config
	if err = yaml.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("yaml.Decode: %w", err)
	}

	return &config, nil
}
