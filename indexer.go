package bitsy

import (
	"context"
	"github.com/viant/bitsy/indexer"
	"github.com/viant/cloudless/data/processor/adapter/gcp"
)

func HandleEvent(ctx context.Context, event gcp.GSEvent) error {
	return indexer.HandleEvent(ctx,event)
}


