package indexer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/viant/afs"
	"github.com/viant/cloudless/data/processor/adapter/gcp"
)

func HandleEvent(ctx context.Context, event gcp.GSEvent) error {
	fs := afs.New()
	service, err := Singleton(ctx, ConfigKey)
	if err != nil {
		fmt.Printf("could not init service: %v\n", err)
		return nil
	}
	request, err := event.NewRequest(ctx, fs, &service.config.Config)
	if err != nil {
		fmt.Printf("build request error: %v\n", err)
		return nil
	}
	report := service.Index(ctx, request)
	jReport, _ := json.Marshal(report)
	fmt.Printf("%s\n", jReport)
	return nil
}
