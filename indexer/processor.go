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
	"math"
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
	intFields := make([]Ints, len(p.IndexingFields))
	floatFields := make([]Floats, len(p.IndexingFields))
	textFields := make([]Texts, len(p.IndexingFields))
	boolFields := make([]Bools, len(p.IndexingFields))
	events := bytes.Split(data, []byte{'\n'})
	err := p.indexEvents(events, textFields, intFields, boolFields, floatFields)
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
	if err = p.logIndexedFloats(floatFields, multiLogger); err != nil {
		return err
	}
	return nil
}

func (p Processor) Post(ctx context.Context, reporter processor.Reporter) error {
	logger := ctx.Value(destination.DataMultiLoggerKey).(*destination.MultiLogger)
	return logger.Close()
}

func (p *Processor) indexEvents(events [][]byte, textFields []Texts, intValues []Ints, boolFields []Bools, floatFields []Floats) error {
	for _, eventLine := range events {

		event := NewEvent(p.Rule)
		if err := gojay.Unmarshal(eventLine, event); err != nil {
			return processor.NewDataCorruption(fmt.Sprintf("invalid json %v, %s", err, eventLine))
		}
		if event.BatchID == math.MinInt64 {
			return processor.NewDataCorruption(fmt.Sprintf("missing value for '%v' batch field: %s", p.BatchField, eventLine))
		}
		if event.Sequence == math.MinInt64 {
			return processor.NewDataCorruption(fmt.Sprintf("missing value for '%v' seq field: %s", p.SequenceField, eventLine))
		}
		for _, field := range event.IndexingFields {
			rawValue, ok := event.values[field.Name]
			if !ok || len(rawValue) == 0 {
				continue
			}

			isRepeated := bytes.HasPrefix(rawValue, []byte("["))
			switch strings.ToLower(field.Type) {
			case config.TypeString:
				if len(textFields[field.Index]) == 0 {
					textFields[field.Index] = make(Texts, 1000)
				}
				if err := p.decodeAndIndexText(isRepeated, rawValue, textFields[field.Index], event); err != nil {
					return err
				}
			case config.TypeInt:
				if len(intValues[field.Index]) == 0 {
					intValues[field.Index] = make(Ints, 10000)
				}
				err := p.decodeAndIndexInt(isRepeated, rawValue, intValues[field.Index], event)
				if err != nil {
					return err
				}

			case config.TypeFloat:
				if len(floatFields[field.Index]) == 0 {
					floatFields[field.Index] = make(Floats, 1000)
				}
				err := p.decodeAndIndexFloat(isRepeated, rawValue, floatFields[field.Index], event)
				if err != nil {
					return err
				}

			case config.TypeBool:
				if len(boolFields[field.Index]) == 0 {
					boolFields[field.Index] = make(Bools, 2)
				}
				err := p.decodeAndIndexBool(isRepeated, rawValue, boolFields[field.Index], event)
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("unsupported type: %s", field.Type)

			}

		}

	}
	return nil
}

func (p *Processor) decodeAndIndexFloat(isRepeated bool, rawValue []byte, values Floats, event *Event) error {
	if isRepeated {
		floats := dec.Floats{Callback: func(value float64) {
			p.indexFloatValues(values, value, event)
		}}
		if err := gojay.Unmarshal(rawValue, &floats); err != nil {
			return fmt.Errorf("failed due to json unmarshall %s,%w", rawValue, err)
		}
	} else {
		value := 0.0
		if err := gojay.Unmarshal(rawValue, &value); err != nil {
			return err
		}
		p.indexFloatValues(values, value, event)
	}
	return nil
}

func (p *Processor) decodeAndIndexBool(isRepeated bool, rawValue []byte, values Bools, event *Event) error {
	if isRepeated {
		bools := dec.Bools{Callback: func(value bool) {
			p.indexBoolValues(values, value, event)
		}}
		if err := gojay.Unmarshal(rawValue, &bools); err != nil {
			return err
		}
	} else {
		value := false
		if err := gojay.Unmarshal(rawValue, &value); err != nil {
			return err
		}
		p.indexBoolValues(values, value, event)
	}
	return nil
}

func (p *Processor) decodeAndIndexInt(isRepeated bool, rawValue []byte, values Ints, event *Event) error {
	if isRepeated {
		ints := dec.Ints{Callback: func(value int) {
			p.indexIntValues(values, value, event)
		}}
		ints.IsQuoted = bytes.Contains(rawValue, []byte(`"`))
		if err := gojay.Unmarshal(rawValue, &ints); err != nil {
			return fmt.Errorf("failed due to json unmarshall %s,%w", rawValue, err)
		}
	} else {
		value := 0
		if err := gojay.Unmarshal(rawValue, &value); err != nil {
			return err
		}
		p.indexIntValues(values, value, event)
	}
	return nil
}

func (p *Processor) decodeAndIndexText(isRepeated bool, rawValue []byte, values Texts, event *Event) error {
	if isRepeated {
		strings := dec.Strings{Callback: func(value string) {
			p.indexTextValues(values, value, event)
		}}
		if err := gojay.Unmarshal(rawValue, &strings); err != nil {
			return err
		}
	} else {
		value := ""
		if err := gojay.Unmarshal(rawValue, &value); err != nil {
			return err
		}
		p.indexTextValues(values, value, event)
	}
	return nil
}

func (p *Processor) logIndexedTexts(texts []Texts, multiLogger *destination.MultiLogger) error {
	for _, field := range p.IndexingFields {
		values := texts[field.Index]
		if len(values) == 0 {
			continue
		}
		key := p.Rule.Dest.TextPrefix + p.Rule.Dest.TableRoot + field.Name
		logger, err := multiLogger.Get(key)
		if err != nil {
			return err
		}
		p.logText(logger, values)
	}
	return nil
}

func (p *Processor) logIndexedNumerics(ints []Ints, multiLogger *destination.MultiLogger) error {
	for _, field := range p.IndexingFields {
		values := ints[field.Index]
		if len(values) == 0 {
			continue
		}
		key := p.Rule.Dest.IntPrefix + p.Rule.Dest.TableRoot + field.Name
		logger, err := multiLogger.Get(key)
		if err != nil {
			return err
		}
		p.logInt(logger, values)
	}
	return nil
}

func (p *Processor) logIndexedFloats(floats []Floats, multiLogger *destination.MultiLogger) error {
	for _, field := range p.IndexingFields {
		values := floats[field.Index]
		if len(values) == 0 {
			continue
		}
		key := p.Rule.Dest.FloatPrefix + p.Rule.Dest.TableRoot + field.Name
		logger, err := multiLogger.Get(key)
		if err != nil {
			return err
		}
		p.logFloat(logger, values)
	}
	return nil
}

func (p Processor) logIndexedBools(bools []Bools, multiLogger *destination.MultiLogger) error {
	for _, field := range p.IndexingFields {
		values := bools[field.Index]
		if len(values) == 0 {
			continue
		}
		key := p.Rule.Dest.BooleanPrefix + p.Rule.Dest.TableRoot + field.Name
		logger, err := multiLogger.Get(key)
		if err != nil {
			return err
		}
		p.logBool(logger, values)
	}
	return nil
}

func (p *Processor) indexTextValues(values Texts, actual string, event *Event) {
	if _, ok := values[actual]; !ok {
		values[actual] = &Text{
			Base: Base{
				Event:  event,
				Events: 0,
			},
			Value: actual,
		}
	}
	textValue := values[actual]
	textValue.Events = textValue.Events | oneBit<<event.Sequence
}

func (p *Processor) indexIntValues(values Ints, actual int, event *Event) {
	if _, ok := values[actual]; !ok {
		values[actual] = &Int{
			Base: Base{
				Event:  event,
				Events: 0,
			},
			Value: actual,
		}
	}
	intValue := values[actual]
	intValue.Events = intValue.Events | oneBit<<(event.Sequence)
}

func (p *Processor) indexFloatValues(values Floats, actual float64, event *Event) {
	if _, ok := values[actual]; !ok {
		values[actual] = &Float{
			Base: Base{
				Event:  event,
				Events: 0,
			},
			Value: actual,
		}
	}
	intValue := values[actual]
	intValue.Events = intValue.Events | oneBit<<(event.Sequence)
}

func (p Processor) indexBoolValues(boolValues Bools, actual bool, event *Event) {
	if _, ok := boolValues[actual]; !ok {
		boolValues[actual] = &Bool{
			Base: Base{
				Event:  event,
				Events: 0,
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

func (p *Processor) logInt(logger *log.Logger, values Ints) {
	for _, value := range values {
		message := p.msgProvider.NewMessage()
		p.logBase(message, &value.Base)
		message.PutInt("value", value.Value)
		logger.Log(message)
		message.Free()
	}
}

func (p *Processor) logFloat(logger *log.Logger, values Floats) {
	for _, value := range values {
		message := p.msgProvider.NewMessage()
		p.logBase(message, &value.Base)
		message.PutFloat("value", value.Value)
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

//NewProcessor creates a new processor
func NewProcessor(rule *config.Rule, concurrency int) *Processor {
	return &Processor{Rule: rule,
		msgProvider: msg.NewProvider(16*1024, concurrency),
	}
}
