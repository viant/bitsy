package indexer

import (
	"github.com/francoispqt/gojay"
	"github.com/viant/bitsy/config"
	"github.com/viant/bitsy/safe"
)

type Event struct {
	*config.Rule
	BatchID   int
	Sequence  int
	Timestamp string
	Partition string
	values    map[string][]byte
	keys      int
}

func (e *Event) UnmarshalJSONObject(dec *gojay.Decoder, key string) (err error) {

	switch key {
	case e.Rule.TimeField:
		return dec.String(&e.Timestamp)
	case e.Rule.PartitionField:
		return dec.String(&e.Partition)
	case e.Rule.BatchField:
		if e.Rule.AllowQuotedNumbers {
			e.BatchID, err = safe.DecodeInt(dec)
			return err
		}
		return dec.Int(&e.BatchID)
	case e.Rule.SequenceField:
		if e.Rule.AllowQuotedNumbers {
			e.Sequence, err = safe.DecodeInt(dec)
			return err
		}
		return dec.Int(&e.Sequence)
	default:
		if _, ok := e.values[key]; ok {
			raw := gojay.EmbeddedJSON{}
			dec.EmbeddedJSON(&raw)
			e.values[key] = raw
		}
	}
	return nil
}




func (e *Event) NKeys() int {
	return e.keys
}


func NewEvent(rule *config.Rule) *Event {
	result := &Event{
		Rule:   rule,
		values: make(map[string][]byte),
		keys:   2 + len(rule.IndexingFields),
	}
	for _, field := range rule.IndexingFields {
		result.values[field.Name] = make([]byte, 0)
	}
	if rule.PartitionField != "" {
		result.keys++
	}
	if rule.TimeField != "" {
		result.keys++
	}
	return result
}