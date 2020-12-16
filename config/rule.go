package config

import (
	"github.com/viant/afs/matcher"
	"github.com/viant/cloudless/data/processor"
	"github.com/viant/cloudless/resource"
)

type Field struct {
	Name string
	Type string
}

type Destination struct {
	URL           string
	TableRoot     string
	TextPrefix    string
	NumericPrefix string
	FloatPrefix   string
	URIKeyName    string
	BooleanPrefix string
}

type Rule struct {
	processor.Config
	TimeField      string
	BatchField     string
	SequenceField  string
	IndexingFields []Field
	Dest           Destination
	Source         matcher.Basic
}

func (r Rule) IsText(field string) bool {
	for _, item := range r.IndexingFields {
		if item.Name == field {
			return item.Type == "string"
		}
	}
	return false
}

type Rules struct {
	BaseURL   string
	CheckInMs int
	Rules     []*Rule
	*resource.Tracker
}
