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
		var cErr error
		service, cErr = NewService(ctx, location)
		if cErr != nil {
			err = fmt.Errorf("failed to create config from env.%v: %v, %w", location, os.Getenv(location), cErr)
			return
		}
	})
	if err != nil {
		runOnce = &sync.Once{}
	}
	return service, err
}


func NewService(ctx context.Context, location string) (*Service, error) {
	fs := afs.New()
	cfg, cErr := config.NewConfigFromEnv(ctx, location)
	if cErr != nil {
		fmt.Errorf("failed to create config from env.%v: %v, %w", location, os.Getenv(location), cErr)
		return nil , err
	}
	return  New(cfg, fs),nil
}


func NewServiceV1(cfg *config.Config,fs afs.Service) (*Service, error) {
	return  New(cfg, fs),nil
}




