package cmd

import (
	"github.com/jessevdk/go-flags"
	"log"
)

//RunClient run client
func RunClient(Version string, args []string) int {
	options := &Options{}
	_, err := flags.ParseArgs(options, args)
	if err != nil {
		log.Fatal(err)
	}
	if options.Validate {
		err :=  validate(options)
		if err !=nil {
			log.Printf("invalid rule %s, %v",options.RuleURL,err)
			return 1
		}
		log.Printf("rule is valid ")
		return 0
	}
	run(options)
	return 0
}

func run(options *Options) {

}
