package isodisplay

import (
	"github.com/tschroed/trafficlight"
	"log"
)

const (
	THRESHOLD_AMBER = 50
	THRESHOLD_RED   = 75
)

type TrafficLightSink struct {
	ch   chan Signal
	done chan bool
	tl   trafficlight.TrafficLight
}

func (s *TrafficLightSink) Input() chan<- Signal {
	return s.ch
}

func (s *TrafficLightSink) Close() error {
	close(s.ch)
	<-s.done
	return nil
}

func NewTrafficLightSink(tl trafficlight.TrafficLight) *TrafficLightSink {
	p := func(s *TrafficLightSink) {
		var oc uint8
		oc = trafficlight.GREEN
		s.tl.Set(oc)
		log.Print("Starting Sink loop...")
		for {
			select {
			case sig, more := <-s.ch:
				if !more {
					s.done <- true
					log.Print("Signal channel closed, exiting Sink loop.")
					return
				}

				log.Printf("Received Signal: %v\n", sig)
				var nc uint8
				nc = trafficlight.GREEN
				if sig.RelativeValue > THRESHOLD_AMBER {
					nc = trafficlight.AMBER
				}
				if sig.RelativeValue > THRESHOLD_RED {
					nc = trafficlight.RED
				}
				if nc != oc {
					s.tl.Set(nc)
					oc = nc
				}
			}
		}
	}
	s := &TrafficLightSink{
		ch:   make(chan Signal, 10),
		done: make(chan bool, 1),
		tl:   tl,
	}
	go p(s)
	return s
}
