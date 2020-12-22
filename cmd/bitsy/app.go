package main

import (
	_ "github.com/viant/afsc/s3"
	"github.com/viant/bitsy/cmd"
	_ "net/http/pprof"
)

var Version string

func main() {
	args := []string{
		"app", "-r=/Users/ppoudyal/go/src/github.com/viant/bitsy/e2e/regression/cases/02_multibatch/rule/rule.yaml", "-s=/Users/ppoudyal/go/src/github.com/viant/bitsy/e2e/regression/cases/02_multibatch/data/trigger/input.json",
	}
	cmd.RunClient(Version, args)

}
