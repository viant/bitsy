package config

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/viant/afs"
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

func NewConfigFromURL(ctx context.Context, fs afs.Service, URL string) (*Config, error) {
	config := &Config{}
	reader, err := fs.OpenURL(ctx, URL)
	if err != nil {
		return nil, fmt.Errorf("failed to open config: %v, due to %w", URL, err)
	}
	defer reader.Close()
	err = json.NewDecoder(reader).Decode(config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config: %v, due to %w", URL, err)
	}
	config.Init()
	config.Config.Init(ctx,fs)
	return config, config.Config.Validate()
}