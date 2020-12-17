package indexer

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/viant/bitsy/config"
	"github.com/viant/cloudless/data/processor"
	"testing"
)

func TestService_Index(t *testing.T) {
	var useCases = []struct {
		description string
		config *config.Config
		request *processor.Request
		expect interface{}
	}{


	}



	for _, useCase := range useCases{
		ctx := context.Background()
		srv := New(useCase.config)
		response := srv.Index(ctx, useCase.request)
		assert.NotNil(t, response, useCase.description)

	}
}
