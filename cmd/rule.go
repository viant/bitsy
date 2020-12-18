package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/viant/afs"
	"github.com/viant/afs/file"
	"github.com/viant/afs/matcher"
	"github.com/viant/afs/url"
	"github.com/viant/bitsy/config"
	"github.com/viant/toolbox"
	"gopkg.in/yaml.v2"
	"path"
	"strconv"
)

func ruleToMap(rule *config.Rule) map[string]interface{} {
	ruleMap := map[string]interface{}{}
	toolbox.DefaultConverter.AssignConverted(&ruleMap, rule)
	compactedMap := map[string]interface{}{}
	toolbox.CopyMap(ruleMap, compactedMap, toolbox.OmitEmptyMapWriter)
	return compactedMap
}

func reportRule(rule *config.Rule) {
	ruleMap := ruleToMap(rule)
	var ruleYAML, err = yaml.Marshal(ruleMap)
	if err == nil {
		fmt.Printf("==== USING RULE ===\n%s===== END ====\n", ruleYAML)
	}
}




func buildRule(options *Options) error {
	fs := afs.New()
	ctx := context.Background()
	options.RuleURL = defaultRuleURL
	basePath, _ := url.Split(options.SourceURL, file.Scheme)
	_, prefix := url.Split(basePath, file.Scheme)
	suffix := path.Ext(options.SourceURL)
	rule := &config.Rule{
		Dest: config.Destination{
			URL: options.DestinationURL,
		},
		BatchField:    options.BatchField,
		SequenceField: options.SequenceField,


		IndexingFields: decodeIndexingFields(options.IndexingFields),
		When: matcher.Basic{
			Prefix: prefix,
			Suffix: suffix,
		},
	}

	if options.SourceURL != "" && len(rule.IndexingFields) ==0 {
		err := autodetectIndexingFields(ctx, options, fs, rule)
		if err != nil {
			return err
		}
	}

	rule.Init()
	data, _ := yaml.Marshal(rule)
	return fs.Upload(ctx, options.RuleURL, file.DefaultFileOsMode, bytes.NewReader(data))
}


func autodetectIndexingFields(ctx context.Context, options *Options, fs afs.Service,  rule *config.Rule) error {
	data, err := fs.DownloadWithURL(ctx, options.SourceURL)
	if err != nil {
		return fmt.Errorf("failed to get source file: %v, %w", options.SourceURL, err)
	}
	aMap := map[string]interface{}{}
	if err := json.Unmarshal(data, &aMap);err != nil {
		return err
	}
	var indexingFields = make(map[string]string)
	for key := range aMap {
		if aMap[key] == nil {
			continue
		}
		switch key {
		case options.SequenceField, options.BatchField, options.TimeField:
			continue
		}
		switch val := aMap[key].(type) {
		case bool:
			indexingFields[key] = config.TypeBool
		case float64:
			indexingFields[key] = config.TypeFloat
		case int64, int, uint, uint64:
			indexingFields[key] = config.TypeInt
		case string:
			if _, err := strconv.Atoi(val); err == nil {
				aMap[key] = config.TypeInt
				rule.AllowQuotedNumbers = true
				continue
			}
			indexingFields[key] = config.TypeString
		}
	}
	rule.IndexingFields = decodeIndexingFields(indexingFields)
	return nil
}



func decodeIndexingFields(indexingFields map[string]string) []config.Field {
	fields := make([]config.Field, 0)
	if len(indexingFields) == 0 {
		return fields
	}
	for k, v := range indexingFields {
		fields = append(fields, config.Field{
			Name: k,
			Type: v,
		})
	}
	return fields
}
