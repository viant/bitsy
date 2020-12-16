package indexer

import (
	"bytes"
	"context"
	"fmt"
	"github.com/francoispqt/gojay"
	"github.com/viant/bitsy/config"
	"github.com/viant/bitsy/indexer/dec"
	"github.com/viant/cloudless/data/processor"
	"github.com/viant/cloudless/data/processor/destination"
	"github.com/viant/tapper/log"
	"github.com/viant/tapper/msg"
	"strings"
)

//Processor represent bitset indexer
type Processor struct {
	*config.Rule
	msgProvider *msg.Provider
}

var oneBit = int64(1)

func (p Processor) Pre(ctx context.Context, reporter processor.Reporter) (context.Context, error) {
	return destination.NewDataMultiLogger(ctx, p.Rule.Dest.URIKeyName, reporter)
}

//Process process data unit (upto 64 rows)
func (p *Processor) Process(ctx context.Context, data []byte, reporter processor.Reporter) error {
	intFields := make(map[string]Ints)
	textFields := make(map[string]Texts)
	boolFields := make(map[string]Bools)
	events := bytes.Split(data, []byte{'\n'})
	err := p.indexEvents(events, textFields, intFields, boolFields)
	if err != nil {
		return err
	}
	multiLogger := ctx.Value(destination.DataMultiLoggerKey).(*destination.MultiLogger)
	if err = p.logIndexedNumerics(intFields, multiLogger); err != nil {
		return err
	}
	if err = p.logIndexedTexts(textFields, multiLogger); err != nil {
		return err
	}
	if err = p.logIndexedBools(boolFields, multiLogger); err != nil {
		return err
	}
	return nil
}

func (p Processor) Post(ctx context.Context, reporter processor.Reporter) error {
	logger := ctx.Value(destination.DataMultiLoggerKey).(*destination.MultiLogger)
	return logger.Close()
}

func (p *Processor) indexEvents(events [][]byte, textFields map[string]Texts, numericFields map[string]Ints, boolFields map[string]Bools) error {
	for _, eventLine := range events {

		event := NewEvent(p.Rule)
		if err := gojay.Unmarshal(eventLine, event); err != nil {
			return processor.NewDataCorruption(fmt.Sprintf("invalid json %v, %s", err, eventLine))
		}

		for _, field := range event.IndexingFields {
			rawValue, ok := event.values[field.Name]
			if !ok {
				continue
			}
			isRepeated := bytes.HasPrefix(rawValue, []byte("["))
			switch strings.ToLower(field.Type) {
			case "string":
				if err := p.decodeAndIndexText(isRepeated, rawValue, textFields, field, event);err != nil {
					return err
				}
			case "int":

				err := p.decodeAndIndexInt(isRepeated, rawValue, numericFields, field, event)
				if err != nil {
					return err
				}

			case "float":

			case "bool":
				err := p.decodeAndIndexBool(isRepeated, rawValue, boolFields, field, event)
				if err != nil {
					return err
				}

			}

		}

		//for field, value := range record {
		//	if value == nil {
		//		continue
		//	}
		//	switch actual := value.(type) {
		//	case bool:
		//		p.indexBoolValues(boolFields, field, actual, timestamp, batchId, seqId)
		//	case string:
		//
		//		if intVal, err := strconv.Atoi(actual); err == nil && !p.Rule.IsText(field) {
		//			p.indexIntValues(numericFields, field, intVal, timestamp, batchId, seqId)
		//		} else {
		//			p.indexTextValues(textFields, field, actual, timestamp, batchId, seqId)
		//		}
		//	case float64:
		//		p.indexIntValues(numericFields, field, int(actual), timestamp, batchId, seqId)
		//	case []float64:
		//		for _, item := range actual {
		//			p.indexIntValues(numericFields, field, int(item), timestamp, batchId, seqId)
		//		}
		//	case []string:
		//		for _, item := range actual {
		//			p.indexTextValues(textFields, field, item, timestamp, batchId, seqId)
		//		}
		//	case []interface{}:
		//		for _, item := range actual {
		//			switch actualItem := item.(type) {
		//			case string:
		//				if intVal, err := strconv.Atoi(actualItem); err == nil && !p.Rule.IsText(field) {
		//					p.indexIntValues(numericFields, field, intVal, timestamp, batchId, seqId)
		//				} else {
		//					p.indexTextValues(textFields, field, actualItem, timestamp, batchId, seqId)
		//				}
		//			case float64:
		//				p.indexIntValues(numericFields, field, int(actualItem), timestamp, batchId, seqId)
		//			}
		//		}
		//	default:
		//		return fmt.Errorf("unsupported type %T", value)
		//	}

	}
	return nil
}

func (p *Processor) decodeAndIndexBool(isRepeated bool, rawValue []byte, boolFields map[string]Bools, field config.Field, event *Event) error {
	if isRepeated {
		bools := dec.Bools{Items: make([]bool, 0)}
		if err := gojay.Unmarshal(rawValue, &bools); err != nil {
			return err
		}
		for _, value := range bools.Items {
			p.indexBoolValues(boolFields, field.Name, value, event)
		}
	} else {
		value := false
		if err := gojay.Unmarshal(rawValue, &value); err != nil {
			return err
		}
		p.indexBoolValues(boolFields, field.Name, value, event)
	}
	return nil
}

func (p *Processor) decodeAndIndexInt(isRepeated bool, rawValue []byte, numericFields map[string]Ints, field config.Field, event *Event) error {
	if isRepeated {
		ints := dec.Ints{Items: make([]int, 0), }
		ints.IsQuoted = bytes.Contains(rawValue,[]byte (`"`))
		if err := gojay.Unmarshal(rawValue, &ints); err != nil {
			return fmt.Errorf("failed due to json unmarshall %s,%w",rawValue,err)
		}
		for _, value := range ints.Items {
			p.indexIntValues(numericFields, field.Name, value, event)
		}
	} else {
		value := 0
		if err := gojay.Unmarshal(rawValue, &value); err != nil {
			return err
		}
		p.indexIntValues(numericFields, field.Name, value, event)
	}
	return nil
}

func (p *Processor) decodeAndIndexText(isRepeated bool, rawValue []byte, textFields map[string]Texts, field config.Field, event *Event) error {
	if isRepeated {
		strings := dec.Strings{Items: make([]string, 0)}
		if err := gojay.Unmarshal(rawValue, &strings); err != nil {
			return err
		}
		for _, value := range strings.Items {
			p.indexTextValues(textFields, field.Name, value, event)
		}
	} else {
		value := ""
		if err := gojay.Unmarshal(rawValue, &value); err != nil {
			return err
		}
		p.indexTextValues(textFields, field.Name, value, event)
	}
	return nil
}

func (p *Processor) logIndexedTexts(textFields map[string]Texts, multiLogger *destination.MultiLogger) error {
	for field, values := range textFields {
		key := p.Rule.Dest.TextPrefix + p.Rule.Dest.TableRoot + field
		logger, err := multiLogger.Get(key)
		if err != nil {
			return err
		}
		p.logText(logger, values)
	}
	return nil
}

func (p *Processor) logIndexedNumerics(numericFields map[string]Ints, multiLogger *destination.MultiLogger) error {
	for field, values := range numericFields {
		key := p.Rule.Dest.NumericPrefix + p.Rule.Dest.TableRoot + field
		logger, err := multiLogger.Get(key)
		if err != nil {
			return err
		}
		p.logNumeric(logger, values)
	}
	return nil
}

func (p Processor) logIndexedBools(boolFields map[string]Bools, multiLogger *destination.MultiLogger) error {

	for field, values := range boolFields {
		key := p.Rule.Dest.BooleanPrefix + p.Rule.Dest.TableRoot + field
		logger, err := multiLogger.Get(key)
		if err != nil {
			return err
		}
		p.logBool(logger, values)
	}
	return nil
}

func (p *Processor) indexTextValues(textFields map[string]Texts, key string, actual string, event *Event) {
	if _, ok := textFields[key]; !ok {
		textFields[key] = Texts{}
	}
	textValues := textFields[key]
	if _, ok := textValues[actual]; !ok {
		textValues[actual] = &Text{
			Base: Base{
				Event: event,
				Events:    0,
			},
			Value: actual,
		}
	}
	textValue := textValues[actual]
	textValue.Events = textValue.Events | oneBit<<event.Sequence
}

func (p *Processor) indexIntValues(intFields map[string]Ints, key string, actual int, event *Event) {
	if _, ok := intFields[key]; !ok {
		intFields[key] = Ints{}
	}
	intValues := intFields[key]
	if _, ok := intValues[actual]; !ok {
		intValues[actual] = &Int{
			Base: Base{
				Event: event,
				Events:    0,
			},
			Value: actual,
		}
	}
	intValue := intValues[actual]
	intValue.Events = intValue.Events | oneBit<<(event.Sequence)
}

func (p Processor) indexBoolValues(boolFields map[string]Bools, key string, actual bool, event *Event) {
	if _, ok := boolFields[key]; !ok {
		boolFields[key] = Bools{}
	}
	boolValues := boolFields[key]
	if _, ok := boolValues[actual]; !ok {
		boolValues[actual] = &Bool{
			Base: Base{
				Event: event,
				Events:    0,
			},
			Value: actual,
		}
	}
	boolValue := boolValues[actual]
	boolValue.Events = boolValue.Events | oneBit<<event.Sequence

}

func (p *Processor) logBase(message *msg.Message, value *Base) {
	message.PutNonEmptyString(p.Rule.TimeField, value.Timestamp)
	message.PutInt(p.Rule.BatchField, value.BatchID)
	message.PutInt("events", int(value.Events))
}

func (p *Processor) logNumeric(logger *log.Logger, values Ints) {
	for _, value := range values {
		message := p.msgProvider.NewMessage()
		p.logBase(message, &value.Base)
		message.PutInt("value", value.Value)
		logger.Log(message)
		message.Free()
	}
}

func (p *Processor) logText(logger *log.Logger, values Texts) {
	for _, value := range values {
		message := p.msgProvider.NewMessage()
		p.logBase(message, &value.Base)
		message.PutString("value", value.Value)
		logger.Log(message)
		message.Free()
	}
}

func (p Processor) logBool(logger *log.Logger, values Bools) {
	for _, value := range values {
		message := p.msgProvider.NewMessage()
		p.logBase(message, &value.Base)
		message.PutBool("value", value.Value)
		logger.Log(message)
		message.Free()
	}

}

//New creates a new processor
func New(rule *config.Rule) *Processor {
	return &Processor{Rule: rule,
		msgProvider: msg.NewProvider(16*1024, rule.Concurrency),
	}
}
