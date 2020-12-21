package main

import (
	_ "github.com/viant/afsc/s3"
	"github.com/viant/bitsy/cmd"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

var Version string

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	cmd.RunClient(Version, os.Args[1:])
}

