package isodisplay

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"
)

const (
	fetchInterval = 10 * time.Second
)

type EmissionsSource struct {
	ch      chan Signal
	fetcher Fetcher
	stop    chan bool
	done    chan bool
}

// Define the inner data structure
type EmissionsData struct {
	NaturalGas  float64
	Oil         float64
	Wood        float64
	Total       float64
	Refuse      float64
	LandfillGas float64
	BeginDateMs int64
	BeginDate   time.Time
}

// Define the outer structure to match the JSON
type emissionsResponse struct {
	Data      []EmissionsData
	Namespace string
}

// ByTime implements sort.Interface for []EmissionsData based on
// the BeginDateMs field.
type ByTime []EmissionsData

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return a[i].BeginDateMs < a[j].BeginDateMs }

func parseEmissionsData(raw []byte) ([]EmissionsData, error) {
	var data []emissionsResponse
	err := json.Unmarshal(raw, &data)
	if err != nil {
		return nil, err
	}
	// log.Printf("parsed: %v (len: %v, namespace: %v)", data, len(data), data[0].Namespace)
	if len(data) != 1 {
		return nil, fmt.Errorf("Don't know how to handle this many responses: %v", len(data))
	}
	sort.Sort(ByTime(data[0].Data))
	return data[0].Data, nil
}

func newEmissionsSource(f Fetcher, interval time.Duration) *EmissionsSource {
	if f == nil {
		log.Print("Received nil fetcher")
		return nil
	}
	s := &EmissionsSource{
		ch:      make(chan Signal, 10),
		fetcher: f,
		stop:    make(chan bool, 1),
		done:    make(chan bool, 1),
	}
	go func() {
		log.Print("Starting Source loop...")
		for {
			select {
			case <-s.stop:
				log.Print("Closing and exiting Source loop")
				close(s.ch)
				s.done <- true
				return
			case <-time.After(interval):
				r, err := f.RawData()
				if err != nil {
					log.Printf("Error fetching data: %v", err)
					continue
				}
				e, err := parseEmissionsData(r)
				if err != nil {
					log.Printf("Error parsing data: %v", err)
					continue
				}
				if len(e) < 1 {
					log.Printf("Emissions data short! %v", e)
					log.Printf("Raw data: %v", string(r))
					continue
				}
				// Cap at 100
				var val int8
				if e[len(e)-1].Total < 100 {
					val = int8(e[len(e)-1].Total)
				} else {
					val = 100
				}
				s.ch <- Signal{
					Name: "TotalEmissions",
					// Right now 0 - 100T/min is a pretty good scale so just grab the raw value directly.
					RelativeValue: int8(val),
					RawValue:      e[len(e)-1].Total,
					RawUnit:       "Metric Tons/min",
				}
			}
		}
	}()
	return s
}

func NewEmissionsSource() *EmissionsSource {
	return newEmissionsSource(NewEmissionsFetcher(), fetchInterval)
}

func (s *EmissionsSource) Output() <-chan Signal {
	return s.ch
}

func (s *EmissionsSource) Close() error {
	s.stop <- true
	<-s.done
	return nil
}
