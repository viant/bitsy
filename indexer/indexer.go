package indexer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/viant/afs"
	"github.com/viant/afs/option"
	"github.com/viant/cloudless/data/processor/adapter/gcp"
)

func HandleEvent(ctx context.Context, event gcp.GSEvent) error {
	fs := afs.New()
	service, err := NewService(ctx, ConfigKey)
	if err != nil {
		fmt.Printf("could not init service: %v\n", err)
		return nil
	}
	request, err := event.NewRequest(ctx, fs, &service.config.Config.Config)
	if err != nil {
		if ok, _ := fs.Exists(ctx, event.URL(), option.NewObjectKind(true)); ! ok {
			fmt.Printf(`{"Status":"noFound", "URL":"%v"}`, event.URL())
			return nil
		}
		return fmt.Errorf("failed to build request: %w", err)
	}
	resp := service.Index(ctx, request)
	JSON, _ := json.Marshal(resp)
	fmt.Printf("%s\n", JSON)
	return nil
}
