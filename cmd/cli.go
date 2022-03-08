package cmd

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
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

func RunApp(Version string,args []string) {
	options := &Options{}
	_, err := flags.ParseArgs(options, args)
	if err != nil {
		log.Fatal(err)
	}
	err = runApp(options)
	if err != nil {
		fmt.Printf("failed to run app: %v\n", err)
	}
	os.Exit(1)

}
