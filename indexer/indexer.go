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
	resp := service.Index(ctx, request)
	JSON, _:= json.Marshal(resp)
	fmt.Printf("%s\n", JSON)
	return nil
}