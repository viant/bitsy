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
	IntPrefix     string
	FloatPrefix   string
	URIKeyName    string
	BooleanPrefix string
}

type Rule struct {
	SourceURL string
	processor.Config
	TimeField          string
	BatchField         string
	SequenceField      string
	PartitionField     string
	IndexingFields     []Field
	AllowQuotedNumbers bool
	Dest               Destination
	Source             matcher.Basic
}
