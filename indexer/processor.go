package indexer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/viant/bitsy/config"
	"github.com/viant/bitsy/index"
	"github.com/viant/cloudless/data/processor"
	"github.com/viant/cloudless/data/processor/destination"
	"github.com/viant/tapper/log"
	"github.com/viant/tapper/msg"
	"github.com/viant/toolbox"
	"strconv"
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
	numericFields := make(map[string]index.Numerics)
	textFields := make(map[string]index.Texts)
	boolFields := make(map[string]index.Bools)
	events := bytes.Split(data, []byte{'\n'})
	err := p.indexEvents(events, textFields, numericFields, boolFields)
	if err != nil {
		return err
	}
	multiLogger := ctx.Value(destination.DataMultiLoggerKey).(*destination.MultiLogger)
	if err = p.logIndexedNumerics(numericFields, multiLogger); err != nil {
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

func (p *Processor) indexEvents(events [][]byte, textFields map[string]index.Texts, numericFields map[string]index.Numerics, boolFields map[string]index.Bools) error {
	for _, event := range events {
		record := make(map[string]interface{})
		err := json.Unmarshal(event, &record)
		if err != nil {
			return processor.NewDataCorruption(fmt.Sprintf("invalid json %v, %s", err, event))
		}
		batchId := toolbox.AsInt(record[p.Rule.BatchField])
		seqId := toolbox.AsInt(record[p.Rule.SequenceField])
		timestamp := toolbox.AsString(record[p.Rule.TimeField])

		for field, value := range record {
			if value == nil {
				continue
			}
			switch actual := value.(type) {
			case bool:
				p.indexBoolValues(boolFields, field, actual, timestamp, batchId, seqId)
			case string:

				if intVal, err := strconv.Atoi(actual); err == nil && !p.Rule.IsText(field) {
					p.indexNumericValues(numericFields, field, intVal, timestamp, batchId, seqId)
				} else {
					p.indexTextValues(textFields, field, actual, timestamp, batchId, seqId)
				}
			case float64:
				p.indexNumericValues(numericFields, field, int(actual), timestamp, batchId, seqId)
			case []float64:
				for _, item := range actual {
					p.indexNumericValues(numericFields, field, int(item), timestamp, batchId, seqId)
				}
			case []string:
				for _, item := range actual {
					p.indexTextValues(textFields, field, item, timestamp, batchId, seqId)
				}
			case []interface{}:
				for _, item := range actual {
					switch actualItem := item.(type) {
					case string:
						if intVal, err := strconv.Atoi(actualItem); err == nil && !p.Rule.IsText(field) {
							p.indexNumericValues(numericFields, field, intVal, timestamp, batchId, seqId)
						} else {
							p.indexTextValues(textFields, field, actualItem, timestamp, batchId, seqId)
						}
					case float64:
						p.indexNumericValues(numericFields, field, int(actualItem), timestamp, batchId, seqId)
					}
				}
			default:
				return fmt.Errorf("unsupported type %T", value)
			}

		}
	}
	return nil
}

func (p *Processor) logIndexedTexts(textFields map[string]index.Texts, multiLogger *destination.MultiLogger) error {
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

func (p *Processor) logIndexedNumerics(numericFields map[string]index.Numerics, multiLogger *destination.MultiLogger) error {
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

func (p Processor) logIndexedBools(boolFields map[string]index.Bools, multiLogger *destination.MultiLogger) error {

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

func (p *Processor) indexTextValues(textFields map[string]index.Texts, key string, actual string, timestamp string, batchId int, seqId int) {
	if _, ok := textFields[key]; !ok {
		textFields[key] = index.Texts{}
	}
	textValues := textFields[key]
	if _, ok := textValues[actual]; !ok {
		textValues[actual] = &index.Text{
			Base: index.Base{
				Timestamp: timestamp,
				BatchID:   batchId,
				Events:    0,
			},
			Value: actual,
		}
	}
	textValue := textValues[actual]
	textValue.Events = textValue.Events | oneBit<<seqId
}

func (p *Processor) indexNumericValues(numericFields map[string]index.Numerics, key string, actual int, timestamp string, batchId int, seqId int) {
	if _, ok := numericFields[key]; !ok {
		numericFields[key] = index.Numerics{}
	}
	numericValues := numericFields[key]
	if _, ok := numericValues[actual]; !ok {
		numericValues[actual] = &index.Numeric{
			Base: index.Base{
				Timestamp: timestamp,
				BatchID:   batchId,
				Events:    0,
			},
			Value: actual,
		}
	}
	numericValue := numericValues[actual]
	numericValue.Events = numericValue.Events | oneBit<<(seqId)
}

func (p Processor) indexBoolValues(boolFields map[string]index.Bools, key string, actual bool, timestamp string, batchId int, seqId int) {
	if _, ok := boolFields[key]; !ok {
		boolFields[key] = index.Bools{}
	}
	boolValues := boolFields[key]
	if _, ok := boolValues[actual]; !ok {
		boolValues[actual] = &index.Bool{
			Base: index.Base{
				Timestamp: timestamp,
				BatchID:   batchId,
				Events:    0,
			},
			Value: actual,
		}
	}
	boolValue := boolValues[actual]
	boolValue.Events = boolValue.Events | oneBit<<seqId

}

func (p *Processor) logBase(message *msg.Message, value *index.Base) {
	message.PutNonEmptyString(p.Rule.TimeField, value.Timestamp)
	message.PutInt(p.Rule.BatchField, value.BatchID)
	message.PutInt("events", int(value.Events))
}

func (p *Processor) logNumeric(logger *log.Logger, values index.Numerics) {
	for _, value := range values {
		message := p.msgProvider.NewMessage()
		p.logBase(message, &value.Base)
		message.PutInt("value", value.Value)
		logger.Log(message)
		message.Free()
	}
}

func (p *Processor) logText(logger *log.Logger, values index.Texts) {
	for _, value := range values {
		message := p.msgProvider.NewMessage()
		p.logBase(message, &value.Base)
		message.PutString("value", value.Value)
		logger.Log(message)
		message.Free()
	}
}

func (p Processor) logBool(logger *log.Logger, values index.Bools) {
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
