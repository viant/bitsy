package indexer

import (
	"context"
	"github.com/viant/bitsy/config"
	"github.com/viant/cloudless/data/processor"
)

type Service struct {
	config *config.Config
}



func (s *Service) Index(ctx context.Context, request *processor.Request) *Reporter {
	return nil
}


func New(cfg *config.Config) *Service {
	return &Service{
		config: cfg,
	}
}