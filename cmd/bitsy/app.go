package main

import (
	_ "github.com/viant/afsc/gs"
	_ "github.com/viant/afsc/s3"
	"github.com/viant/bitsy/cmd"
	_ "net/http/pprof"
	"os"
)

var Version string = "1.0"
func main() {
	args := os.Args
	cmd.RunApp(Version, args)
}
