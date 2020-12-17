package config

import (
	"fmt"
	"github.com/viant/afs/matcher"
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
	SourceURL          string
	TimeField          string
	BatchField         string
	SequenceField      string
	PartitionField     string
	IndexingFields     []Field
	AllowQuotedNumbers bool
	Dest               Destination
	Source             matcher.Basic
}

func (r *Rule) Init() {
	r.Dest.Init()
}

func (d *Destination) Init() {
	if d.URIKeyName == "" {
		d.URIKeyName = "$fragment"
	}
	if d.IntPrefix == "" {
		d.IntPrefix = "num/"
	}
	if d.FloatPrefix == "" {
		d.FloatPrefix = "float/"
	}
	if d.TextPrefix == "" {
		d.TextPrefix = "text/"
	}
	if d.BooleanPrefix == "" {
		d.BooleanPrefix = "bool/"
	}

}

func (r *Rule) Validate() error {
	if len(r.BatchField) == 0 {
		return fmt.Errorf("batch field is missing")
	}
	return nil
}
