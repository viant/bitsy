package indexer

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/viant/afs"
	"github.com/viant/assertly"
	"github.com/viant/bitsy/config"
	"github.com/viant/cloudless/data/processor"
	"github.com/viant/toolbox"
	"io/ioutil"
	"path"
	"strings"
	"testing"
)

func TestService_Index(t *testing.T) {

	parent := toolbox.CallerDirectory(3)

	var useCases = []struct {
		description  string
		config       *config.Config
		request      *processor.Request
		expect       interface{}
		expectedData map[string]string
	}{
		{
			description: "yaml rule",
			config: &config.Config{
				Rules: config.Rules{
					BaseURL: path.Join(parent, "test_data/index/01_yaml/rules"),
					Config: processor.Config{
						RetryURL:        "mem://localhost/bitsy/retry",
						FailedURL:       "mem://localhost/bitsy/failed",
						CorruptionURL:   "mem://localhost/bitsy/corrupted",
						MaxExecTimeMs:   120000,
						ScannerBufferMB: 2,
						MaxRetries:      2,
					},
				},
			},
			request: &processor.Request{
				SourceURL: "mem://localhost/case001/data.json",
				ReadCloser: ioutil.NopCloser(strings.NewReader(`{"id": 1, "name": "Adam", "country": "us", "city_id":1, "batch_id":1, "seq":0, "tstamp":"2020-11-01 01:01:01"}
{"id": 2, "name": "Kent", "country":"us","batch_id":1, "city_id":1, "seq":1, "tstamp":"2020-11-01 01:01:01"}
{"id": 3, "name": "Adam", "country":"nep", "batch_id":1,"city_id":2, "seq":2, "tstamp":"2020-11-01 01:01:01"}`)),
			},
			expect: `{
	"Batched": 1,
	"CorruptionURL": "mem://localhost/bitsy/corrupted/case001/data.json",
	"DestinationURL": "mem://localhost/index/case01/$fragment/data.json",
	"Loaded": 3,
	"Processed": 1,
	"RetryURL": "mem://localhost/bitsy/retry/case001/data-retry01.json",
	"SourceURL": "mem://localhost/case001/data.json",
	"Status": "ok"
}`,
			expectedData: map[string]string{
				"mem://localhost/index/case01/int/myTable_city_id/data.json": `[{"@indexBy@": "value"},
{"batch_id":1, "value":1, "events":3 },
{"batch_id":1, "value":2, "events":4 }]`,
			},
		},
	}

	fs := afs.New()
	for _, useCase := range useCases {

		ctx := context.Background()
		srv := New(useCase.config, fs)

		response := srv.Index(ctx, useCase.request)
		assert.NotNil(t, response, useCase.description)
		if !assertly.AssertValues(t, useCase.expect, response, useCase.description) {
			toolbox.DumpIndent(response, true)
		}

		if len(useCase.expectedData) == 0 {
			continue
		}
		for URL, expect := range useCase.expectedData {
			actual, err := fs.DownloadWithURL(ctx, URL)
			if !assert.Nil(t, err, useCase.description+" / "+URL) {
				continue
			}
			expectedIndex := []interface{}{}
			err = json.Unmarshal([]byte(expect), &expectedIndex)
			if !assert.Nil(t, err, useCase.description+" / "+URL) {
				continue
			}
			if !assertly.AssertValues(t, expectedIndex, string(actual), useCase.description+" / "+URL) {
				toolbox.DumpIndent(expectedIndex, true)
			}

		}
	}

}
