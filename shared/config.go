package shared

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL     string `json:"DatabaseURL"`
	LogLevel        string `json:"LogLevel"`
	ProducerPort    int    `json:"ProducerPort"`
	ConsumerPort    int    `json:"ConsumerPort"`
	MaxBacklog      int    `json:"MaxBacklog"`
	PrometheusPort  int    `json:"PrometheusPort"`
	ConsumerAddress string `json:"ConsumerAddress"`
}

func LoadConfig() (*Config, error) {
	file, err := os.Open("../shared/config.json")
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("error decoding config file: %v", err)
	}
	return &config, nil
}
