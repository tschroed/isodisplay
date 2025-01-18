package isodisplay

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	emissionsData = `[{"data":[
		{"NaturalGas":29.16,"Oil":0.24,"Wood":6.4,"Total":44.62,"Refuse":8.47,"LandfillGas":0.35,"BeginDateMs":1737198711000,"BeginDate":"2025-01-18T06:11:51.000-05:00"},
		{"NaturalGas":26.69,"Oil":0.03,"Wood":6.42,"Total":41.67,"Refuse":8.18,"LandfillGas":0.35,"BeginDateMs":1737186246000,"BeginDate":"2025-01-18T02:44:06.000-05:00"},
		{"NaturalGas":34.01,"Oil":0.95,"Wood":6.46,"Total":50.26,"Refuse":8.49,"LandfillGas":0.35,"BeginDateMs":1737176922000,"BeginDate":"2025-01-18T00:08:42.000-05:00"},
		{"NaturalGas":35.41,"Oil":1.91,"Wood":6.44,"Total":52.26,"Coal":0.05,"Refuse":8.1,"LandfillGas":0.35,"BeginDateMs":1737223920000,"BeginDate":"2025-01-18T13:12:00.000-05:00"},
		{"NaturalGas":33.48,"Oil":0.78,"Wood":6.31,"Total":49.37,"Coal":0.03,"Refuse":8.42,"LandfillGas":0.35,"BeginDateMs":1737207448000,"BeginDate":"2025-01-18T08:37:28.000-05:00"}
	], "namespace": "_nstmp_"}]`
)

func TestParseEmissionsData(t *testing.T) {
	// Slightly out of order to test that result is sorted.
	raw := []byte(emissionsData)
	data, err := parseEmissionsData(raw)
	log.Printf("Got data: %v", data)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(data) != 5 {
		t.Errorf("Unexpected number of entries: %v", len(data))
	}
	expectedTotals := []float64{50.26, 41.67, 44.62, 49.37, 52.26}
	for i, tot := range expectedTotals {
		if tot != data[i].Total {
			t.Errorf("Result %v, unexpected total. Wanted %v, got %v (%v)", i, tot, data[i].Total, data[i])
		}
	}

	raw = []byte(`[{}]`)
	data, err = parseEmissionsData(raw)
	log.Printf("Got data: %v", data)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(data) != 0 {
		t.Errorf("Unexpected number of entries: %v", len(data))
	}

	raw = []byte(`"Garbage"`)
	data, err = parseEmissionsData(raw)
	log.Printf("Got data: %v", data)
	if err == nil {
		t.Errorf("Expected error: %v", err)
	}

	// Too many data sets.
	raw = []byte(`[{"data":[], "namespace": "_nstmp_"},{"data":[], "namespace": "foo"}]`)
	data, err = parseEmissionsData(raw)
	log.Printf("Got data: %v", data)
	if err == nil {
		t.Errorf("Expected error: %v", err)
	}
}

func TestNewEmissionsSource(t *testing.T) {
	log.Print("Creating EmissionsSource...")
	s := NewEmissionsSource()
	log.Print("Closing EmissionsSource...")
	s.Close()
}

func Test_newEmissionsSource(t *testing.T) {
	s := newEmissionsSource(nil, 12345)
	if s != nil {
		t.Errorf("Expected nil EmissionsSource, got %v", s)
		s.Close()
	}
}

func fetcherAndServerForTest() (Fetcher, *httptest.Server) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, emissionsData)
	}))
	url := func(*httpFetcher) string {
		return s.URL
	}
	now := func() time.Time {
		return time.Now()
	}
	f := newHTTPFetcher(url, now, http.Get, 5*time.Minute)
	return f, s
}

func TestEmissionsSource(t *testing.T) {
	ftch, srvr := fetcherAndServerForTest()
	defer srvr.Close()
	s := newEmissionsSource(ftch, 1*time.Second)
	if s == nil {
		t.Errorf("nil EmissionsSource)}")
	}
	sig := <-s.Output()
	log.Printf("Signal: %v", sig)
	s.Close()
}
