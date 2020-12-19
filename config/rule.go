package config

import (
	"fmt"
	"github.com/viant/afs/file"
	"github.com/viant/afs/matcher"
	"github.com/viant/afs/url"
	"path"
	"strings"
	"time"
)

type Field struct {
	Name  string
	Type  string
	Index int
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
	Disabled           bool
	SourceURL          string
	TimeField          string
	BatchField         string
	SequenceField      string
	PartitionField     string
	IndexingFields     []Field
	fields             map[string]*Field
	AllowQuotedNumbers bool
	Dest               Destination
	When               matcher.Basic
}

func (r *Rule) Fields() map[string]*Field {
	return r.fields
}

//HasMatch returns true if URL matches prefix or suffix
func (r *Rule) HasMatch(URL string) bool {
	location := url.Path(URL)
	parent, name := path.Split(location)
	match := r.When.Match(parent, file.NewInfo(name, 0, 0644, time.Now(), false))
	return match
}

func (r *Rule) Init() {
	r.Dest.Init()
	if len(r.IndexingFields) == 0 {
		return
	}
	r.fields = map[string]*Field{}
	for i, item := range r.IndexingFields {
		r.IndexingFields[i].Index = i
		r.fields[item.Name] = &r.IndexingFields[i]
	}
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

func (d *Destination) Validate() error {
	if d.URL == "" {
		return fmt.Errorf("destination URL was empty")
	}
	if !strings.Contains(d.URL, d.URIKeyName) {
		return fmt.Errorf("destionaionURL %v doesn't contain %v", d.URL, d.URIKeyName)
	}

	return nil
}

func (r *Rule) Validate() error {
	err := r.Dest.Validate()
	if err != nil {
		return err
	}
	if r.BatchField == "" {
		return fmt.Errorf("batchfield was empty")
	}
	if r.SequenceField == "" {
		return fmt.Errorf("sequencefield was empty")
	}
	if r.When.Prefix == "" {
		return fmt.Errorf("when.Prefix was empty")
	}
	if len(r.IndexingFields) == 0 {
		return fmt.Errorf("indexingFields was empty")
	}

	for _, field := range r.IndexingFields {
		switch strings.ToLower(field.Type) {
		case TypeBool, TypeFloat, TypeInt, TypeString:
		default:
			return fmt.Errorf("unsupported data type: '%s', for field: %s, supported: %v", field.Type, field.Name, []string{
				TypeBool, TypeFloat, TypeInt, TypeString,
			})
		}
	}
	return nil
}
