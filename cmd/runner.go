package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/viant/afs"
	"github.com/viant/afs/file"
	"github.com/viant/afs/url"
	"github.com/viant/bitsy/config"
	"github.com/viant/bitsy/indexer"
	"github.com/viant/cloudless/data/processor"
	"math"
	"os"
)

const (
	configKey = "CONFIG"
)

func run(options *Options) error {

	fs := afs.New()
	ctx := context.Background()
	parentURL, _ := url.Split(options.RuleURL, file.Scheme)
	cfg := config.Config{
		Rules: config.Rules{
			BaseURL: parentURL,
			Config: processor.Config{
				Concurrency:         options.Concurrency,
				DestinationURL:      options.DestinationURL,
				DestinationCodec:    "",
				RetryURL:            "file:///tmp/bitsy/retry",
				FailedURL:           "file:///tmp/bitsy/failed",
				CorruptionURL:      "file:///tmp/bitsy/corrupted",
				MaxExecTimeMs:       math.MaxInt64,
			},
		},
	}
	JSON, _ := json.Marshal(cfg)
	os.Setenv(configKey, string(JSON))
	srv, err := indexer.Singleton(ctx, configKey)
	if err != nil {
		return err
	}
	reader, err := fs.OpenURL(ctx, options.SourceURL)
	if err != nil {
		return err
	}

	defer reader.Close()
	response := srv.Index(ctx, &processor.Request{
		SourceURL: options.SourceURL,
		ReadCloser: reader,
	})
	JSON, _ = json.Marshal(response)
	fmt.Printf("%s\n", response)
	return nil
}

/*
	srv := New(useCase.config, fs)

 */
