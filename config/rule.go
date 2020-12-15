package config

import (
	"github.com/viant/afs/matcher"
	"github.com/viant/bitsy/resource"
	"github.com/viant/cloudless/data/processor"
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

type Rules struct {
	BaseURL   string
	CheckInMs int
	Rules     []*Rule
	meta      *resource.Tracker
}
