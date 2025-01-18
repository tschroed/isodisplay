package isodisplay

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

func TestNewStdoutSink(t *testing.T) {
	log.Println("Creating Sink")
	s := NewStdoutSink()
	if s == nil {
		t.Errorf("s is nil")
	}
	if s.out != os.Stdout {
		t.Errorf("s.out not stdout: %v", s.out)
	}
	log.Println("Closing Sink")
	if err := s.Close(); err != nil {
		t.Errorf("Error on Close(): %v", err)
	}
}

func TestStdoutSinkInput(t *testing.T) {
	var buf bytes.Buffer
	log.Println("Creating Sink")
	s := newStdoutSink(&buf)
	s.Input() <- Signal{
		Name:          "foo",
		RelativeValue: 42,
		RawValue:      1337,
		RawUnit:       "bar",
	}
	log.Println("Closing Sink")
	if err := s.Close(); err != nil {
		t.Errorf("Error on Close(): %v", err)
	}
	log.Printf("buf: %v", buf.String())
	if !strings.Contains(buf.String(), "{foo 42 1337 bar}") {
		t.Errorf("Output did not contain expected text")
	}
}
