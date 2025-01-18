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
	dec232024us = 1734985560123456
)

const (
	OK = iota
	ERR_500
	ERR_FETCH
	ERR_NIL
	OK_ALT
)

func TestEmissionsURL(t *testing.T) {
	ts := time.UnixMicro(dec232024us)
	want := "https://www.iso-ne.com/ws/wsclient?_nstmp_formDate=1734985560123&_nstmp_startDate=12/23/2024&_nstmp_endDate=12/23/2024&_nstmp_twodays=false&_nstmp_requestType=emissions"
	got := emissionsURL(ts)
	if want != got {
		t.Errorf("Unexpected emissions URL: got %s, want %s", got, want)
	}
}

func TestHTTPFetcher(t *testing.T) {
	expected := "dummy data"
	expected2 := "dummy data too"
	mode := OK
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch mode {
		case ERR_500:
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "This is an error message")
		case OK_ALT:
			fmt.Fprintf(w, expected2)
		default:
			fmt.Fprintf(w, expected)
		}
	}))
	defer svr.Close()
	url := func(*httpFetcher) string {
		return svr.URL
	}
	var n int64
	n = dec232024us
	now := func() time.Time {
		return time.UnixMicro(n)
	}
	get := func(url string) (*http.Response, error) {
		switch mode {
		case OK, ERR_500, OK_ALT:
			return http.Get(url)
		case ERR_FETCH:
			return nil, fmt.Errorf("Fake error fetching: %s", url)
		case ERR_NIL:
			return nil, nil
		default:
			panic(fmt.Sprintf("Unknown case: %v", mode))
		}
	}
	f := newHTTPFetcher(url, now, get, 5*time.Minute)
	// Normal fetch
	d, err := f.RawData()
	log.Printf("Got raw data: %v", d)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if d != expected {
		t.Errorf("Unexpected data, got \"%s\", expected \"%s\"", d, expected)
	}
	d, err = f.RawData()

	// Set error mode, but cache should mask it
	mode = ERR_500
	log.Printf("Got raw data: %v", d)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if d != expected {
		t.Errorf("Unexpected data, got \"%s\", expected \"%s\"", d, expected)
	}

	// Now expire the cache, so new fetch should propagate 500
	n += 60000000 * 6 // Cache is expired 6 minutes later
	d, err = f.RawData()
	log.Printf("Got raw data: %v", d)
	if err == nil {
		t.Errorf("Expected error, received none")
	}

	// Now error on the fetch itsel
	mode = ERR_FETCH
	n += 60000000 * 6 // Cache is expired 6 minutes later
	d, err = f.RawData()
	log.Printf("Got raw data: %v", d)
	if err == nil {
		t.Errorf("Expected error, received none")
	}

	// Now return a nil response
	mode = ERR_NIL
	n += 60000000 * 6 // Cache is expired 6 minutes later
	d, err = f.RawData()
	log.Printf("Got raw data: %v", d)
	if err == nil {
		t.Errorf("Expected error, received none")
	}

	// And flush the cache letting a new body be fetched
	mode = OK_ALT
	f.FlushCache()
	d, err = f.RawData()
	log.Printf("Got raw data: %v", d)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if d != expected2 {
		t.Errorf("Unexpected data, got \"%s\", expected \"%s\"", d, expected)
	}
}

func TestNewEmissionsFetcher(t *testing.T) {
	NewEmissionsFetcher()
}
