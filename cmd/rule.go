package cmd

import (
	"fmt"
	"github.com/viant/bitsy/config"
	"github.com/viant/toolbox"
	"gopkg.in/yaml.v2"
)

func ruleToMap(rule *config.Rule) map[string]interface{} {
	ruleMap := map[string]interface{}{}
	toolbox.DefaultConverter.AssignConverted(&ruleMap, rule)
	compactedMap := map[string]interface{}{}
	toolbox.CopyMap(ruleMap, compactedMap, toolbox.OmitEmptyMapWriter)
	return compactedMap
}

func reportRule(rule *config.Rule) {
	ruleMap := ruleToMap(rule)
	var ruleYAML, err = yaml.Marshal(ruleMap)
	if err == nil {
		fmt.Printf("==== USING RULE ===\n%s===== END ====\n", ruleYAML)
	}
}
