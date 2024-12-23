package isodisplay

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

const (
	dec232024ms = 1734985560123456
)

func TestEmissionsURL(t *testing.T) {
	ts := time.UnixMicro(dec232024ms)
	want := "https://www.iso-ne.com/ws/wsclient?_nstmp_formDate=1734985560123&_nstmp_startDate=12/23/2024&_nstmp_endDate=12/23/2024&_nstmp_twodays=false&_nstmp_requestType=emissions"
	got := emissionsURL(ts)
	if want != got {
		t.Errorf("Unexpected emissions URL: got %s, want %s", got, want)
	}
}

func TestHTTPFetcher(t *testing.T) {
	u := "https://foo/bar"
	url := func(*httpFetcher) string {
		return u
	}
	now := func() time.Time {
		return time.UnixMicro(dec232024ms)
	}
	// TODO(tschroed): Update this to use an httptest.Server and Client from that.
	get := func(url string) (*http.Response, error) {
		if url != u {
			return nil, fmt.Errorf("Unexpected URL: %s", url)
		}
		return nil, nil
	}
	f := newHTTPFetcher(url, now, get, 5*time.Minute)
	_, err := f.RawData()
	if err == nil {
		t.Errorf("Expected an error, got none.")
	}
}
