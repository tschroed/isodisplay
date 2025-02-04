package isodisplay

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Fetcher interface {
	RawData() ([]byte, error)
}

type httpFetcher struct {
	url    func(*httpFetcher) string
	now    func() time.Time
	get    func(string) (*http.Response, error)
	ttl    time.Duration
	cache  []byte
	expiry time.Time
}

func newHTTPFetcher(url func(*httpFetcher) string, now func() time.Time, get func(string) (*http.Response, error), ttl time.Duration) *httpFetcher {
	return &httpFetcher{
		url: url,
		now: now,
		get: get,
		ttl: ttl,
	}
}

func (f *httpFetcher) FlushCache() {
	f.expiry = time.UnixMicro(0)
}

func (f *httpFetcher) RawData() ([]byte, error) {
	now := f.now()
	if now.UnixNano() > f.expiry.UnixNano() {
		url := f.url(f)
		log.Printf("Fetching %s", url)
		resp, err := f.get(url)
		if err != nil {
			msg := fmt.Errorf("Error fetching %s: %w", url, err)
			log.Printf("%v", msg)
			return nil, msg
		}
		if resp == nil {
			msg := fmt.Errorf("Got nil response")
			log.Printf("%v", msg)
			return nil, msg
		}
		defer resp.Body.Close()
		log.Printf("Response %d bytes (code: %d)", resp.ContentLength, resp.StatusCode)
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("Unexpected response fetching %s: %s (%d)", url, resp.Status, resp.StatusCode)
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Error reading body %s: %w", url, err)
		}
		log.Printf("Read %d bytes", len(b))
		f.cache = b
		f.expiry = f.now().Add(f.ttl)
	}
	return f.cache, nil
}

func emissionsURL(t time.Time) string {
	tmpl := "https://www.iso-ne.com/ws/wsclient?_nstmp_formDate=%d&_nstmp_startDate=%s&_nstmp_endDate=%s&_nstmp_twodays=false&_nstmp_requestType=emissions"
	today := t.Format("01/02/2006")
	return fmt.Sprintf(tmpl, t.UnixMilli(), today, today)
}

func NewEmissionsFetcher() *httpFetcher {
	url := func(f *httpFetcher) string {
		return emissionsURL(f.now())
	}
	return newHTTPFetcher(url, time.Now, http.Get, 10*time.Minute)
}
