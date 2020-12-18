package cmd

import (
	"context"
	"fmt"
	"github.com/viant/afs"
	"github.com/viant/afs/file"
	"github.com/viant/bitsy/config"
	"path"
)

func validate(options *Options) error{
	fs := afs.New()
	data,err := fs.DownloadWithURL(context.Background(),options.RuleURL,file.DefaultFileOsMode)
	if err != nil {
		return fmt.Errorf("failed to load %s, %w", options.RuleURL,err)
	}
	rule, err := config.LoadRule(data,path.Ext(options.RuleURL))
	if err != nil {
		return err
	}
	rule.Init()
	return rule.Validate()
}
