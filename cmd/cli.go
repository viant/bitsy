package cmd

import (
	"github.com/jessevdk/go-flags"
	"log"
)

const (
	defaultRuleURL = "mem://localhost/bitsy/rules/rule.yaml"
)

//RunClient run client
func RunClient(Version string, args []string) int {
	options := &Options{}
	_, err := flags.ParseArgs(options, args)
	if err != nil {
		log.Fatal(err)
	}

	if options.RuleURL == "" {
		buildRule(options)
	}

	if options.Validate {
		err := validate(options)
		if err != nil {
			log.Printf("invalid rule %s, %v", options.RuleURL, err)
			return 1
		}
		log.Printf("Rule is VALID")
		return 0
	}
	if err := run(options); err != nil {
		log.Print(err)
		return 1
	}
	return 0
}
