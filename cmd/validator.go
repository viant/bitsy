package cmd

import (
	"context"
	"fmt"
	"github.com/viant/afs"
	"github.com/viant/afs/file"
	"github.com/viant/bitsy/config"
	"path"
)

func validate(options *Options) error {
	fs := afs.New()
	rule, err := loadRule(options, fs)
	if err != nil {
		return err
	}
	err = rule.Validate()
	if err != nil {
		return err
	}
	reportRule(rule)
	return nil
}

func loadRule(options *Options, fs afs.Service) (*config.Rule, error) {
	data, err := fs.DownloadWithURL(context.Background(), options.RuleURL, file.DefaultFileOsMode)
	if err != nil {
		return nil, fmt.Errorf("failed to load %s, %w", options.RuleURL, err)
	}
	rule, err := config.LoadRule(data, path.Ext(options.RuleURL))
	if err != nil {
		return nil, err
	}
	rule.Init()
	return rule, nil
}
