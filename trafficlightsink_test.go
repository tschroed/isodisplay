package isodisplay

import (
	"testing"

	"github.com/tschroed/trafficlight"
)

type mockTrafficLight struct {
	ch    chan bool
	state uint8
}

func (m *mockTrafficLight) Set(word uint8) error {
	m.state = word
	m.ch <- true
	return nil
}

func TestTrafficLightSink(t *testing.T) {
	mockTL := &mockTrafficLight{
		ch: make(chan bool, 1),
	}
	sink := NewTrafficLightSink(mockTL)

	tests := []struct {
		signal        Signal
		expectedState uint8
	}{
		{Signal{RelativeValue: 30}, trafficlight.GREEN},
		{Signal{RelativeValue: 60}, trafficlight.AMBER},
		{Signal{RelativeValue: 80}, trafficlight.RED},
	}

	for _, test := range tests {
		sink.Input() <- test.signal
		<-mockTL.ch
		if mockTL.state != test.expectedState {
			t.Errorf("Expected state %v, got %v", test.expectedState, mockTL.state)
		}
	}

	sink.Close()
}
