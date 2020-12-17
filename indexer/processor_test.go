package indexer

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/afs"
	"github.com/viant/afs/url"
	"github.com/viant/assertly"
	"github.com/viant/bitsy/config"
	"github.com/viant/cloudless/data/processor"
	"testing"
	"time"
)

func Test_Process(t *testing.T) {

	var useCases = []struct {
		description string
		config.Rule
		config.Rules
		input       string
		expectedURL string
		expected    map[string]string
		hasError    bool
	}{

		{
			description: "multi rows index",
			Rules: config.Rules{
				Config: processor.Config{
					DestinationURL: "mem://localhost/data/$fragment/data.json",
				},
			},

			Rule: config.Rule{
				Dest: config.Destination{
					TableRoot:     "test_",
					TextPrefix:    "text/",
					IntPrefix:     "num/",
					FloatPrefix:   "float/",
					URIKeyName:    "$fragment",
					BooleanPrefix: "bool/",
				},
				BatchField:    "batch_id",
				SequenceField: "seq",
				TimeField:     "tstamp",
				IndexingFields: []config.Field{
					{
						Name: "name",
						Type: "string",
					},
					{
						Name: "country",
						Type: "string",
					},
					{
						Name: "city_id",
						Type: "int",
					},
				},
			},
			expectedURL: "mem://localhost/data/",
			input: `{"id": 1, "name": "Adam", "country": "us", "city_id":1, "batch_id":1, "seq":0, "tstamp":"2020-11-01 01:01:01"}
{"id": 2, "name": "Kent", "country":"us","batch_id":1, "city_id":1, "seq":1, "tstamp":"2020-11-01 01:01:01"}
{"id": 3, "name": "Adam", "country":"nep", "batch_id":1,"city_id":2, "seq":2, "tstamp":"2020-11-01 01:01:01"}`,
			expected: map[string]string{
				"text/test_name/data.json": `{"@indexBy@": "value"}
{"batch_id":1, "value":"Adam", "events":5 }
{"batch_id":1, "value":"Kent", "events":2 }`,
				"text/test_country/data.json": `{"@indexBy@": "value"}
{"batch_id":1, "value":"us", "events":3 }
{"batch_id":1, "value":"nep", "events":4 }
`,
				"num/test_city_id/data.json": `{"@indexBy@": "value"}
{"batch_id":1, "value":1, "events":3 }
{"batch_id":1, "value":2, "events":4 }
`,
			},
		},
		{
			description: "repeated rows index",
			Rules: config.Rules{
				Config: processor.Config{
					DestinationURL: "mem://localhost/case2/$fragment/data.json",
				},
			},
			Rule: config.Rule{
				Dest: config.Destination{
					URL:           "",
					TableRoot:     "test_",
					TextPrefix:    "text/",
					IntPrefix:     "num/",
					FloatPrefix:   "float/",
					URIKeyName:    "$fragment",
					BooleanPrefix: "bool/",
				},
				BatchField:    "batch_id",
				SequenceField: "seq",
				TimeField:     "tstamp",
				IndexingFields: []config.Field{
					{
						Name: "segments",
						Type: "int",
					},
				},
				AllowQuotedNumbers: true,
			},
			expectedURL: "mem://localhost/case2/",
			input: `{"id": 1, "segments": ["1","10","100"],"batch_id":1, "seq":0, "tstamp":"2020-11-01 01:01:01"}
{"id": 2, "segments": [1,20],"batch_id":1, "seq":1, "tstamp":"2020-11-01 01:01:01"}
{"id": 3, "segments": [1,10,100], "batch_id":1,"seq":2, "tstamp":"2020-11-01 01:01:01"}`,
			expected: map[string]string{
				"num/test_segments/data.json": `{"@indexBy@": "value"}
{"batch_id":1, "value":1, "events":7 }
{"batch_id":1, "value":10, "events":5 }
{"batch_id":1, "value":20, "events":2 }
{"batch_id":1, "value":100, "events":5 }
`,
			},
		},
		{
			description: "boolean rows index",
			Rules: config.Rules{
				Config: processor.Config{
					DestinationURL: "mem://localhost/case2/$fragment/data.json",
				},
			},
			Rule: config.Rule{
				Dest: config.Destination{
					URL:           "",
					TableRoot:     "test_",
					TextPrefix:    "text/",
					IntPrefix:     "num/",
					FloatPrefix:   "float/",
					URIKeyName:    "$fragment",
					BooleanPrefix: "bool/",
				},
				BatchField:    "batch_id",
				SequenceField: "seq",
				TimeField:     "tstamp",
				IndexingFields: []config.Field{
					{
						Name: "is_pmp",
						Type: "bool",
					},
				},
			},
			expectedURL: "mem://localhost/case2/",
			input: `{"id": 1, "is_pmp": true,"batch_id":1, "seq":0, "tstamp":"2020-11-01 01:01:01"}
{"id": 2, "is_pmp": true,"batch_id":1, "seq":1, "tstamp":"2020-11-01 01:01:01"}
{"id": 3, "is_pmp": false, "batch_id":1,"seq":2, "tstamp":"2020-11-01 01:01:01"}`,
			expected: map[string]string{
				"bool/test_is_pmp/data.json": `{"@indexBy@": "value"}
{"batch_id":1, "value":true, "events":3 }
{"batch_id":1, "value":false, "events":4 }
`,
			},
		},
		{
			description: "float rows index",
			Rules: config.Rules{
				Config: processor.Config{
					DestinationURL: "mem://localhost/case4/$fragment/data.json",
				},
			},
			Rule: config.Rule{
				Dest: config.Destination{
					URL:           "",
					TableRoot:     "test_",
					TextPrefix:    "text/",
					IntPrefix:     "num/",
					FloatPrefix:   "float/",
					URIKeyName:    "$fragment",
					BooleanPrefix: "bool/",
				},
				BatchField:    "batch_id",
				SequenceField: "seq",
				TimeField:     "tstamp",
				IndexingFields: []config.Field{
					{
						Name: "cp",
						Type: "float",
					},
				},
			},
			expectedURL: "mem://localhost/case4/",
			input: `{"id": 1, "cp": 1.2E10-5,"batch_id":1, "seq":0, "tstamp":"2020-11-01 01:01:01"}
{"id": 2, "cp": null,"batch_id":1, "seq":1, "tstamp":"2020-11-01 01:01:01"}
{"id": 3, "cp": 0.1, "batch_id":1,"seq":2, "tstamp":"2020-11-01 01:01:01"}`,
			expected: map[string]string{
				"float/test_cp/data.json": `{"@indexBy@": "value"}
{"batch_id":1, "value":0.000012, "events":1 }
{"batch_id":1, "value":0.1, "events":4 }
{"batch_id":1, "value":0.0, "events":2 }
`,
			},
		},
	}

	fs := afs.New()

	for _, useCase := range useCases {
		ctx := context.Background()
		proc := NewProcessor(&useCase.Rule, 10)

		reporter := processor.NewReporter()
		reporter.BaseResponse().DestinationURL = useCase.Rules.ExpandDestinationURL(time.Now())

		ctx, err := proc.Pre(ctx, reporter)
		if !assert.Nil(t, err, useCase.description) {
			continue
		}
		err = proc.Process(ctx, []byte(useCase.input), reporter)
		if useCase.hasError {
			assert.NotNil(t, err, useCase.description)
			continue
		}
		if !assert.Nil(t, err, useCase.description) {
			continue
		}

		err = proc.Post(ctx, reporter)
		if !assert.Nil(t, err, useCase.description) {
			continue
		}

		for URI, content := range useCase.expected {
			URL := url.Join(useCase.expectedURL, URI)
			actual, err := fs.DownloadWithURL(ctx, URL)
			if !assert.Nil(t, err, useCase.description+" / "+URL) {
				continue
			}

			if !assertly.AssertValues(t, content, string(actual), useCase.description+" / "+URL) {
				fmt.Printf("%s\n", actual)
			}
		}

	}

}
