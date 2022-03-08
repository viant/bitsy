package config

import (
	"context"
	"github.com/viant/afs"
	"github.com/viant/cloudless/data/processor"
	"github.com/viant/cloudless/data/processor/subscriber/gcp"
	"github.com/viant/cloudless/resource"
	"log"
	"path"
	"sync"
	"time"
)

type Rules struct {
	gcp.Config
	BaseURL          string
	CheckFrequencyMs int
	Indexes          []*Rule
	*resource.Tracker
	mux       sync.RWMutex
	PreSorted *bool
}

func (r *Rules) Match(URL string) []*Rule {
	r.mux.RLock()
	defer r.mux.RUnlock()
	var result = make([]*Rule, 0)
	for _, candidate := range r.Indexes {
		if candidate.HasMatch(URL) {
			if candidate.Disabled {
				continue
			}
			result = append(result, candidate)
		}
	}
	return result
}

func (r *Rules) ReloadIfNeeded(ctx context.Context, fs afs.Service) error {
	var rules = make(map[string]*Rule)
	hasChanged := false
	err := r.Notify(ctx, fs, func(URL string, operation resource.Operation) {
		hasChanged = true
		if len(rules) == 0 {
			r.mux.RLock()
			for i, rule := range r.Indexes {
				rules[rule.SourceURL] = r.Indexes[i]
			}
			r.mux.RUnlock()
		}

		switch operation {
		case resource.Added, resource.Modified:
			rule, err := r.loadRule(ctx, URL, fs)
			if err != nil {
				log.Printf("failed to load %v, %v\n", URL, err)
				return
			}
			rule.SourceURL = URL
			rules[rule.SourceURL] = rule

		case resource.Deleted:
			delete(rules, URL)
		}

		return
	})
	if err != nil || !hasChanged {
		return err
	}

	//Convert rules to r.Indexes
	var updatedRules = make([]*Rule, 0)
	for key := range rules {
		updatedRules = append(updatedRules, rules[key])
	}
	r.mux.Lock()
	defer r.mux.Unlock()
	r.Indexes = updatedRules
	return nil
}

func (r *Rules) logCallBackError(err error) {
	log.Printf("error executing reload callback %v\n", err)
}

func (r *Rules) loadRule(ctx context.Context, URL string, fs afs.Service) (*Rule, error) {
	data, err := fs.DownloadWithURL(ctx, URL)
	if err != nil {
		return nil, err
	}
	rule, err := LoadRule(data, path.Ext(URL))
	if err != nil {
		return nil, err
	}
	rule.Init()
	return rule, rule.Validate()
}

func (r *Rules) Init() {
	r.Indexes = make([]*Rule, 0)
	r.Tracker = resource.New(r.BaseURL, time.Duration(r.CheckFrequencyMs)*time.Microsecond)
	if r.Config.MaxExecTimeMs == 0 {
		r.Config.MaxExecTimeMs = 3600000
	}
	if r.Config.ScannerBufferMB == 0 {
		r.ScannerBufferMB = 2
	}
	if r.Config.Concurrency == 0 {
		r.Config.Concurrency = 100
	}
	if r.CheckFrequencyMs == 0 {
		r.CheckFrequencyMs = 60000
	}
}

func (r *Rules) ProcessorConfig(rule *Rule) processor.Config {
	cfg := r.Config.Config
	cfg.DestinationURL = rule.Dest.URL
	cfg.DestinationCodec = rule.Dest.Codec
	cfg.BatchSize = 64
	preSorted := rule.PreSorted
	if r.PreSorted != nil && *r.PreSorted {
		preSorted = true
	}
	if !preSorted {
		cfg.Sort.By = []processor.Field{
			{
				Name:      rule.BatchField,
				IsNumeric: true,
			},
		}
		cfg.Sort.Format = "json"
		cfg.Sort.Batch = true
	}

	return r.Config.Config
}
