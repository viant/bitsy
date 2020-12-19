package main

import (
	_ "github.com/viant/afsc/s3"
	"github.com/viant/bitsy/cmd"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	//	"os"
)

var Version string

func main() {
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()
	//cmd.RunClient("", []string{
	//	"-r=/tmp/bitsy/rule/valid.yaml", "-s=/tmp/bitsy/data.json",
	//})
	cmd.RunClient(Version, os.Args[1:])

	time.Sleep(time.Hour)
}

