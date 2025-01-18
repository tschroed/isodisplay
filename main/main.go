package main

// TODO(trevors): Flags for different sinks.
import (
	"log"

	"github.com/tschroed/isodisplay"
	"github.com/tschroed/trafficlight/lcus"
)

func main() {
	log.Print("Setting up sources and sinks...")
	tl, err := lcus.New("/dev/ttyUSB0")
	if err != nil {
		panic(err)
	}
	src := isodisplay.NewEmissionsSource()
	snk := isodisplay.NewTrafficLightSink(tl)
	log.Print("Starting event loop.")
	for {
		sig := <-src.Output()
		snk.Input() <- sig
	}
}
