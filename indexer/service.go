package indexer

import (
	"context"
	"fmt"
	"github.com/viant/afs"
	"github.com/viant/bitsy/config"
	"github.com/viant/cloudless/data/processor"
)

type Service struct {
	config *config.Config
	fs     afs.Service
}

func (s *Service) Index(ctx context.Context, request *processor.Request) *Reporter {
	reporter := NewReporter()
	err := s.index(ctx, request, reporter)
	if err != nil {
		reporter.BaseResponse().LogError(err)
	}
	return reporter
}

func (s *Service) index(ctx context.Context, request *processor.Request, reporter *Reporter) error {
	err := s.config.ReloadIfNeeded(ctx, s.fs)
	if err != nil {
		return err
	}


	rules := s.config.Match(request.SourceURL)
	switch len(rules) {

	case 0:
		reporter.BaseResponse().Status = StatusNoMatch
		return nil
	case 1:

		cfg := s.config.ProcessorConfig(rules[0])
		proc := NewProcessor(rules[0],s.config.Concurrency)
		srv := processor.New(&cfg, s.fs, proc, func() processor.Reporter {
			return reporter
		})
		srv.Do(ctx, request)
		return nil
	default:
		return fmt.Errorf("too many rules matched %+v", rules)

	}
}

func New(cfg *config.Config, fs afs.Service) *Service {
	cfg.Init()
	return &Service{
		config: cfg,
		fs: fs,
	}
}
