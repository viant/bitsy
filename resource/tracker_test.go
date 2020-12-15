package resource

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/afs"
	"github.com/viant/afs/asset"
	"github.com/viant/afs/file"
	"log"
	"testing"
	"time"
)

func TestTracker_HasChanged(t *testing.T) {

	var useCases = []struct {
		description       string
		baseURL           string
		resources         []*asset.Resource
		modifications     []*asset.Resource
		expectedURL       string
		expectedOperation Operation
		checkFrequency    time.Duration
	}{
		{
			description: "test addition url ",
			baseURL:     "mem://localhost/case1",
			resources: []*asset.Resource{
				asset.NewFile("abc.json", []byte("foo bar"), file.DefaultFileOsMode),
			},
			modifications: []*asset.Resource{
				asset.NewFile("def.json", []byte(" car sar"), file.DefaultFileOsMode),
			},
			expectedURL:       "mem://localhost/case1/abc.json",
			expectedOperation: OperationDeleted,
			checkFrequency:    1 * time.Second,
		},
	}
	ctx := context.Background()
	fs := afs.New()
	for _, useCase := range useCases {
		mgr, err := afs.Manager(useCase.baseURL)
		if err != nil {
			log.Fatal(err)
		}
		err = asset.Create(mgr, useCase.baseURL, useCase.resources)
		if err != nil {
			log.Fatal(err)
		}
		tracker := New(useCase.baseURL, useCase.checkFrequency)
		initialResourcesCount := 0
		fmt.Printf("is next change : %v\n", tracker.nextCheck.IsZero())
		err = tracker.Notify(ctx, fs, func(URL string, operation Operation) {
			initialResourcesCount++
		})
		assert.Nil(t, err, useCase.description)
		assert.Equal(t, len(useCase.resources), initialResourcesCount, useCase.description)

		err = asset.Create(mgr, useCase.baseURL, useCase.modifications)
		if err != nil {
			log.Fatal(err)
		}
		actualURL := ""
		var actualOperation Operation = OperationUndefined
		time.Sleep(2 * time.Second)
		err = tracker.Notify(ctx, fs, func(URL string, operation Operation) {
			actualURL = URL
			actualOperation = operation
		})

		assert.EqualValues(t, useCase.expectedURL, actualURL, useCase.description)
		assert.EqualValues(t, useCase.expectedOperation, actualOperation, useCase.description)

	}
}
