package core

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/vitalvas/junos-streaming-analytics/internal/output"
	"gopkg.in/yaml.v3"
)

type CollectorConfig struct {
	Outputs         map[string]output.Config      `yaml:"outputs" json:"outputs"`
	JTI             map[string]CollectorJTIConfig `yaml:"jti" json:"jti"`
	Workers         int                           `yaml:"workers" json:"workers"`
	SendOutputTimer time.Duration                 `yaml:"send_output_timer" json:"send_output_timer"`
}

type CollectorJTIConfig struct {
	Addr   string   `yaml:"addr" json:"addr"`
	Output []string `yaml:"output" json:"output"`
}

func ParseCollectorConfig(configPath string) (*CollectorConfig, error) {
	config := CollectorConfig{
		Outputs: map[string]output.Config{
			"default": {
				Type: "console",
			},
		},
		JTI: map[string]CollectorJTIConfig{
			"default": {
				Addr: "[::]:21000",
			},
		},
		Workers:         runtime.NumCPU(),
		SendOutputTimer: 10 * time.Second,
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Println("config file not found, using default config")
		return &config, nil
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	switch {
	case strings.HasSuffix(configPath, ".json"):
		if err := json.NewDecoder(file).Decode(&config); err != nil {
			return nil, err
		}

	case strings.HasSuffix(configPath, ".yaml"), strings.HasSuffix(configPath, ".yml"):
		if err := yaml.NewDecoder(file).Decode(&config); err != nil {
			return nil, err
		}

	default:
		return nil, errors.New("unknown config file format")
	}

	return &config, nil
}
