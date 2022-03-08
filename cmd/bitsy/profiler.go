//go:build profiler
// +build profiler

package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

func init() {
	go func() {
		log.Println(http.ListenAndServe(":8081", nil))
	}()
}
