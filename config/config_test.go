package config

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/viant/afs"
	"github.com/viant/toolbox"
	"testing"
)

func TestNewConfigFromURL(t *testing.T) {
	var useCases = []struct {
		description string
		URL string
	} {
		{
			description: "test config",
			URL: "/Users/ppoudyal/go/src/github.com/viant/bitsy/e2e/pubsub/app/config.json",
		},
	}
	for _, useCase := range useCases {
		ctx := context.Background()
		fs := afs.New()
		cfg, err := NewConfigFromURL(ctx,fs,useCase.URL)
		assert.Nil(t, err,useCase.description)
		toolbox.Dump(cfg.Config.Config)
	}
}
