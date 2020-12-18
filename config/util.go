package config

import (
	"encoding/json"
	"github.com/viant/toolbox"
	"gopkg.in/yaml.v2"
)

func LoadRule(data []byte, ext string) (*Rule, error) {
	if ext == "" {
		return nil, nil
	}
	rule := &Rule{}

	switch ext {
	case ".yaml":
		ruleMap := map[string]interface{}{}
		if err := yaml.Unmarshal(data, &ruleMap); err != nil {
			rulesMap := []map[string]interface{}{}
			err = json.Unmarshal(data, &rulesMap)
			if err != nil {
				return nil, err
			}
			err = toolbox.DefaultConverter.AssignConverted(&rule, rulesMap)
			return rule, err
		}
		err := toolbox.DefaultConverter.AssignConverted(&rule, ruleMap)
		return rule, err
	default:
		err := json.Unmarshal(data, &rule)
		return rule, err

	}
}
