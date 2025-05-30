package bitsy

import (
	"context"
	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/viant/bitsy/indexer"
	"github.com/viant/cloudless/data/processor/adapter/gcp"
)

func HandleEvent(ctx context.Context, event gcp.GSEvent) error {
	return indexer.HandleEvent(ctx, event)
}
