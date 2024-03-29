package cmd

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"github.com/viant/afs"
	"github.com/viant/afs/file"
	"github.com/viant/afs/url"
	_ "github.com/viant/afsc/gs"
	_ "github.com/viant/afsc/s3"
	"github.com/viant/bitsy/config"
	"github.com/viant/bitsy/indexer"
	"github.com/viant/cloudless/data/processor"
	"github.com/viant/cloudless/data/processor/subscriber/gcp"
	"github.com/viant/cloudless/ioutil"
	"math"
	"os"
	"time"
)

const (
	configKey = "CONFIG"
)

func run(options *Options) error {
	fs := afs.New()
	rule, err := loadRule(options, fs)
	if err != nil {
		return err
	}
	reportRule(rule)
	ctx := context.Background()
	parentURL, _ := url.Split(options.RuleURL, file.Scheme)
	cfg := config.Config{
		Rules: config.Rules{
			BaseURL: parentURL,
			Config: gcp.Config{
				Config: processor.Config{
					Concurrency:      options.Concurrency,
					DestinationURL:   options.DestinationURL,
					DestinationCodec: "",
					RetryURL:         "file:///tmp/bitsy/retry",
					FailedURL:        "file:///tmp/bitsy/failed",
					CorruptionURL:    "file:///tmp/bitsy/corrupted",
					MaxExecTimeMs:    math.MaxInt32,
				},
			},
		},
	}
	cfg.Init()
	JSON, _ := json.Marshal(cfg)
	os.Setenv(configKey, string(JSON))
	srv, err := indexer.Singleton(ctx, configKey)
	if err != nil {
		return err
	}
	reader, err := ioutil.OpenURL(ctx, fs, options.SourceURL)
	if err != nil {
		return err
	}
	defer reader.Close()
	response := srv.Index(ctx, &processor.Request{
		SourceURL:  options.SourceURL,
		ReadCloser: reader,
		StartTime:  time.Now(),
	})
	JSON, _ = json.Marshal(response)
	fmt.Printf("%s\n", JSON)
	return nil
}

//pubsub entry
func runApp(options *Options) error {
	fs := afs.New()
	ctx := context.Background()
	cfg, err := config.NewConfigFromURL(ctx, fs, options.ConfigURL)
	if err != nil {
		return err
	}
	srv, err := indexer.NewServiceV1(cfg, fs)
	if err != nil {
		return err
	}
	client, err := pubsub.NewClient(context.Background(), cfg.ProjectID)
	if err != nil {
		return err
	}

	proc := processor.NewHandler(srv.Handle)
	procService := processor.New(&cfg.Config.Config, fs, proc, func() processor.Reporter {
		reporter := indexer.NewReporter()

		return processor.NewHandlerReporter(reporter)
	})
	subscriberService, err := gcp.New(&cfg.Config, client, procService, fs)
	if err != nil {
		return err
	}
	return subscriberService.Consume(ctx)
}
