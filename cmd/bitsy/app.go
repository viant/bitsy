package main

import (
	_ "github.com/viant/afsc/s3"
	"github.com/viant/bitsy/cmd"
	_ "net/http/pprof"
	"os"
)

var Version string = "1.0"

func main() {
	//args := []string{
	//	"app", "-r=/Users/ppoudyal/go/src/github.com/viant/bitsy/e2e/regression/cases/02_multibatch/rule/rule.yaml", "-s=/Users/ppoudyal/go/src/github.com/viant/bitsy/e2e/regression/cases/02_multibatch/data/trigger/input.json",
	//}

	args := os.Args
	cmd.RunClient(Version, args)

}
