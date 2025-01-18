package main

import (
	"flag"
	"log"

	"github.com/tschroed/isodisplay"
	"github.com/tschroed/trafficlight/lcus"
)

var dFlag = flag.String("d", "lcus", "Driver to use (stdout or lcus)")

func main() {
	flag.Parse()
	log.Print("Setting up sources and sinks...")
	src := isodisplay.NewEmissionsSource()
	var snk isodisplay.Sink
	switch *dFlag {
	case "lcus":
		tl, err := lcus.New("/dev/ttyUSB0")
		if err != nil {
			panic(err)
		}
		snk = isodisplay.NewTrafficLightSink(tl)
	case "stdout":
		snk = isodisplay.NewStdoutSink()
	default:
		panic("Unknown driver: " + *dFlag)
	}
	log.Print("Starting event loop.")
	for {
		sig := <-src.Output()
		snk.Input() <- sig
	}
}
