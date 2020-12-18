package main

import (
	_ "github.com/viant/afsc/s3"
	"github.com/viant/bitsy/cmd"
	"os"

	//	"os"
)

var Version string

func main() {
	//cmd.RunClient("", []string{
	//	"-r=/tmp/bitsy/rule/valid.yaml", "-s=/tmp/bitsy/data.json",
	//})
	cmd.RunClient(Version, os.Args[1:])
}

