package config

import (
	"context"
	"github.com/viant/afs"
	"github.com/viant/cloudless/resource"
	"log"
	"sync"
	"time"
)

type Rules struct {
	processor.Config
	BaseURL   string
	CheckInMs int
	Indexes   []*Rule
	*resource.Tracker
	mux sync.RWMutex
}



func (r *Rules) Match(URL string) []*Rule {
	r.mux.RLock()
	defer r.mux.RUnlock()
	return nil
}

func (r *Rules) ReloadIfNeeded(ctx context.Context, fs afs.Service) {
	r.mux.Lock()
	defer r.mux.Unlock()


	var rules = make(map[string]*Rule)
	hasChanged := false
	//TODO return bool, error in case notified was called ?
	r.Notify(ctx, fs, func(URL string, operation resource.Operation) {
		hasChanged = true
		if len(rules) == 0 {
			for i, rule := range r.Indexes {
				rules[rule.SourceURL] = r.Indexes[i]
			}
		}


		switch operation {
			case resource.OperationAdded:
				 rule, err := r.loadRule(ctx, URL, fs);
				 if err != nil {
					log.Printf("failed to load %v, %v\n", URL, err)
					continue
				}
				rules[rule.SourceURL]= rule

			case resource.OperationDeleted:

			case resource.OperationModified:

		}

	})
	//Convert rules to r.Indexes

}

func (r *Rules) loadRule(ctx context.Context,  URL string, fs afs.Service) (*Rule, error) {
	return nil, nil
}

func (r *Rules) Init()  {
	r.Indexes = make([]*Rule, 0)
	r.Tracker = resource.New(r.BaseURL, time.Duration(r.CheckInMs) * time.Microsecond )
}
