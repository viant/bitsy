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
	floatFields := make(map[string]Floats)
	textFields := make(map[string]Texts)
	boolFields := make(map[string]Bools)
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

func (p *Processor) indexEvents(events [][]byte, textFields map[string]Texts, numericFields map[string]Ints, boolFields map[string]Bools, floatFields map[string]Floats) error {
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
				if err := p.decodeAndIndexText(isRepeated, rawValue, textFields, field, event); err != nil {
					return err
				}
			case "int":

				err := p.decodeAndIndexInt(isRepeated, rawValue, numericFields, field, event)
				if err != nil {
					return err
				}

			case "float":
				err := p.decodeAndIndexFloat(isRepeated, rawValue, floatFields, field, event)
				if err != nil {
					return err
				}

			case "bool":
				err := p.decodeAndIndexBool(isRepeated, rawValue, boolFields, field, event)
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

func (p *Processor) decodeAndIndexFloat(isRepeated bool, rawValue []byte, floatFields map[string]Floats, field config.Field, event *Event) error {
	if isRepeated {
		floats := dec.Floats{Callback: func(value float64) {
			p.indexFloatValues(floatFields, field.Name, value, event)
		}}
		if err := gojay.Unmarshal(rawValue, &floats); err != nil {
			return fmt.Errorf("failed due to json unmarshall %s,%w", rawValue, err)
		}
	} else {
		value := 0.0
		if err := gojay.Unmarshal(rawValue, &value); err != nil {
			return err
		}
		p.indexFloatValues(floatFields, field.Name, value, event)
	}
	return nil
}

func (p *Processor) decodeAndIndexBool(isRepeated bool, rawValue []byte, boolFields map[string]Bools, field config.Field, event *Event) error {
	if isRepeated {
		bools := dec.Bools{Callback: func(value bool) {
			p.indexBoolValues(boolFields, field.Name, value, event)
		}}
		if err := gojay.Unmarshal(rawValue, &bools); err != nil {
			return err
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
		ints := dec.Ints{Callback: func(value int) {
			p.indexIntValues(numericFields, field.Name, value, event)
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
		p.indexIntValues(numericFields, field.Name, value, event)
	}
	return nil
}

func (p *Processor) decodeAndIndexText(isRepeated bool, rawValue []byte, textFields map[string]Texts, field config.Field, event *Event) error {
	if isRepeated {
		strings := dec.Strings{Callback: func(value string) {
			p.indexTextValues(textFields, field.Name, value, event)
		}}
		if err := gojay.Unmarshal(rawValue, &strings); err != nil {
			return err
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
		key := p.Rule.Dest.IntPrefix + p.Rule.Dest.TableRoot + field
		logger, err := multiLogger.Get(key)
		if err != nil {
			return err
		}
		p.logInt(logger, values)
	}
	return nil
}

func (p *Processor) logIndexedFloats(floatFields map[string]Floats, multiLogger *destination.MultiLogger) error {
	for field, values := range floatFields {
		key := p.Rule.Dest.FloatPrefix + p.Rule.Dest.TableRoot + field
		logger, err := multiLogger.Get(key)
		if err != nil {
			return err
		}
		p.logFloat(logger, values)
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
				Event:  event,
				Events: 0,
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
				Event:  event,
				Events: 0,
			},
			Value: actual,
		}
	}
	intValue := intValues[actual]
	intValue.Events = intValue.Events | oneBit<<(event.Sequence)
}

func (p *Processor) indexFloatValues(floatFields map[string]Floats, key string, actual float64, event *Event) {
	if _, ok := floatFields[key]; !ok {
		floatFields[key] = Floats{}
	}
	floatValues := floatFields[key]
	if _, ok := floatValues[actual]; !ok {
		floatValues[actual] = &Float{
			Base: Base{
				Event:  event,
				Events: 0,
			},
			Value: actual,
		}
	}
	intValue := floatValues[actual]
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
