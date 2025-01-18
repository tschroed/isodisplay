package isodisplay

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// A Sink that outputs the signal to stdout
type StdoutSink struct {
	ch   chan Signal
	done chan bool
	out  io.Writer
}

func (s *StdoutSink) Input() chan<- Signal {
	return s.ch
}

func (s *StdoutSink) Close() error {
	close(s.ch)
	<-s.done
	return nil
}

func newStdoutSink(out io.Writer) *StdoutSink {
	p := func(s *StdoutSink) {
		log.Print("Starting Sink loop...")
		for {
			select {
			case sig, more := <-s.ch:
				if !more {
					log.Print("Signal channel closed, exiting Sink loop.")
					s.done <- true
					return
				}
				fmt.Fprintf(s.out, "[%v] Received Signal: %v\n", time.Now(), sig)
			}
		}
	}
	s := &StdoutSink{
		ch:   make(chan Signal),
		done: make(chan bool),
		out:  out,
	}
	go p(s)
	return s
}

func NewStdoutSink() *StdoutSink {
	return newStdoutSink(os.Stdout)
}
