package config

import (
	"fmt"
	"github.com/viant/afs/matcher"
	"strings"
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

func(d *Destination) Validate() error {
	if d.URL == "" {
		return fmt.Errorf("destination URL was empty")
	}
	if !strings.Contains(d.URL, d.URIKeyName) {
		return fmt.Errorf("destionaionURL %v doesn't contain %v", d.URL, d.URIKeyName)
	}

	return nil
}

func (r *Rule) Validate() error {
	err := r.Validate()
	if err != nil {
		return err
	}
	if r.BatchField == "" {
		return fmt.Errorf("batchfield was empty")
	}
	if r.SequenceField == "" {
		return fmt.Errorf("sequencefield was empty")
	}
	return nil
}
