package config

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"strings"
)

type Config struct {
	Rules
}

//NewConfigFromEnv creates a new config from env
func NewConfigFromEnv(ctx context.Context, key string) (*Config, error) {
	data := os.Getenv(key)
	if data == "" {
		return nil, fmt.Errorf("config data was empty")
	}
	cfg := &Config{}
	err := json.NewDecoder(strings.NewReader(data)).Decode(cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode config: %v", cfg)
	}
	return cfg, nil
}
