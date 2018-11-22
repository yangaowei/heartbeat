package main

//
import (
	"./logs"
	"flag"
)

var (
	port    string
	pattern string //api  cmd
	debug   bool   //api  cmd
	//url     string
)

func initFlag() {
	flag.StringVar(&port, "port", "8002", "server port")
	flag.BoolVar(&debug, "debug", false, "logs pattern")
	flag.Parse()
}

func main() {
	initFlag()
	//
	if debug {
		logs.Log.SetLevel(8)
	} else {
		logs.Log.SetLevel(8)
	}
	logs.Log.Debug("pattern: %s", pattern)

}
