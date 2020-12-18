package indexer

import (
	"context"
	"fmt"
	"github.com/viant/afs"
	"github.com/viant/bitsy/config"
	"os"
	"sync"
)

var service *Service
var err error

var runOnce = &sync.Once{}

func Singleton(ctx context.Context, location string) (*Service, error) {

	runOnce.Do(func() {
		fs := afs.New()
		cfg, cErr := config.NewConfigFromEnv(ctx, location)
		if cErr != nil {
			err = fmt.Errorf("failed to create config from env.%v: %v, %w", location, os.Getenv(location), cErr)
			return
		}
		service = New(cfg, fs)

	})

	if err != nil {
		runOnce = &sync.Once{}
	}
	return service, err
}
