package main

//
import (
	"./api"
	"./logs"
	"./worker"
	"flag"
)

var (
	port      string
	pattern   string //api  cmd
	debug     bool   //api  cmd
	workerNum int
	//url     string
)

func initFlag() {
	flag.StringVar(&pattern, "pattern", "api", "web api")
	flag.BoolVar(&debug, "debug", false, "logs pattern")
	flag.IntVar(&workerNum, "wn", 0, "workerNum")
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
		worker.Run(workerNum)
	} else {
		logs.Log.Debug("Please use the mode in（api,worker）")
	}

}
