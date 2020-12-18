package cmd

import (
	"bytes"
	"context"
	"github.com/jessevdk/go-flags"
	"github.com/viant/afs"
	"github.com/viant/afs/file"
	"github.com/viant/afs/matcher"
	"github.com/viant/afs/url"
	"github.com/viant/bitsy/config"
	"gopkg.in/yaml.v2"
	"log"
	"path"
	"strings"
)

const (
	defaultRuleURL = "mem://localhost/bitsy/rules/rule.yaml"
)

//RunClient run client
func RunClient(Version string, args []string) int {
	options := &Options{}
	_, err := flags.ParseArgs(options, args)
	if err != nil {
		log.Fatal(err)
	}

	if options.RuleURL == "" {
		buildRule(options)
	}

	if options.Validate {
		err := validate(options)
		if err != nil {
			log.Printf("invalid rule %s, %v", options.RuleURL, err)
			return 1
		}
		log.Printf("Rule is VALID")
		return 0
	}
	if err := run(options); err != nil {
		log.Print(err)
		return 1
	}
	return 0
}

func buildRule(options *Options) {
	fs := afs.New()
	options.RuleURL = defaultRuleURL
	basePath, _ := url.Split(options.SourceURL, file.Scheme)
	_, prefix := url.Split(basePath, file.Scheme)
	suffix := path.Ext(options.SourceURL)
	rule := &config.Rule{
		Dest: config.Destination{
			URL: options.DestinationURL,
		},
		BatchField:     options.BatchField,
		SequenceField:  options.SequenceField,
		IndexingFields: parseIndexingFields(options.IndexingFields),
		When: matcher.Basic{
			Prefix: prefix,
			Suffix: suffix,
		},
	}
	rule.Init()
	data, _ := yaml.Marshal(rule)
	fs.Upload(context.Background(), options.RuleURL, file.DefaultFileOsMode, bytes.NewReader(data))
}

func parseIndexingFields(sFields string) []config.Field {
	fields := make([]config.Field, 0)
	data := strings.Split(sFields, ",")
	for _, item := range data {
		nameAndType := strings.Split(item, ":")
		if len(nameAndType) == 2 {
			fields = append(fields, config.Field{
				Name: nameAndType[0],
				Type: nameAndType[1],
			})
		}

	}
	return fields
}
