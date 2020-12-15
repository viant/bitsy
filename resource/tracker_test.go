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
		modifcations      []*asset.Resource
		expectedURL       string
		expectedOperation int
		checkFrequency   time.Duration

	}{
		{
			description:       "test addition url ",
			baseURL:           "mem://localhost/case1",
			resources:        [] *asset.Resource{
				asset.NewFile("abc.txt",[]byte ("foo bar") ,file.DefaultFileOsMode),
			} ,
			modifcations: [] *asset.Resource {
				asset.NewFile("def.txt",[]byte (" car sar"), file.DefaultFileOsMode),
			},
			expectedURL: "mem://localhost/case1/def.txt",
			expectedOperation: 0,
			checkFrequency: 1*time.Second,
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
		tracker := New(useCase.baseURL,useCase.checkFrequency)
		initialResourcesCount := 0
		fmt.Printf("is next change : %v\n",tracker.nextCheck.IsZero())
		err = tracker.Notify(ctx,fs, func(URL string, operation int) {
			initialResourcesCount++
		})
		assert.Nil(t,err,useCase.description)
		assert.Equal(t,len(useCase.resources),initialResourcesCount,useCase.description)
		err = asset.Create(mgr, useCase.baseURL, useCase.modifcations)
		if err != nil {
			log.Fatal(err)
		}
		actualURL := ""
		actualOperation := -1
		time.Sleep(2*time.Second)
		err = tracker.Notify(ctx,fs, func(URL string, operation int) {
			actualURL = URL
			actualOperation = operation
		})

		assert.EqualValues(t,useCase.expectedURL,actualURL,useCase.description)
		assert.EqualValues(t,useCase.expectedOperation,actualOperation,useCase.description)


	}
}
