package main

//
import (
	"./api"
	"./logs"
	"./worker"
	"flag"
)

var (
	port    string
	pattern string //api  cmd
	debug   bool   //api  cmd
	//url     string
)

func initFlag() {
	flag.StringVar(&pattern, "pattern", "api", "web api")
	flag.BoolVar(&debug, "debug", false, "logs pattern")
	flag.Parse()
}

func main() {
	initFlag()

	if debug {
		logs.Log.SetLevel(8)
	} else {
		logs.Log.SetLevel(8)
	}
	logs.Log.Debug("pattern: %s", pattern)
	if pattern == "api" {
		api.Run()
	} else if pattern == "worker" {
		var v string
		logs.Log.Debug(v)
		worker.Run()
	} else {
		logs.Log.Debug("Please use the mode in（api,worker）")
	}

}
