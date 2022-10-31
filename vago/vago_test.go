package vago

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createRequest(url string, referrer string, userAgent string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, url, nil)
	h := req.Header

	if referrer != "" {
		h.Set("Referer", referrer)
	}
	if userAgent != "" {
		h.Set("User-Agent", userAgent)
	}
	return req
}

func fixTime(event *VAEvent) {
	event.Timestamp = time.Date(2020, time.May, 12, 23, 50, 21, 0, time.UTC).UnixMilli()
}

func TestFromRequest(t *testing.T) {
	req := createRequest("https://example.com/test?param=abc", "https://www.google.com/search", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	event := FromRequest(req)
	fixTime(&event)

	assert.Equal(t, "https", event.Scheme)
	assert.Equal(t, "example.com", event.Host)
	assert.Equal(t, "https://www.google.com/search", event.Referrer)
	assert.Equal(t, "Mozilla/5.0 (Windows NT 10.0; Win64; x64)", event.UserAgent)

	json, _ := json.Marshal(event)
	assert.Equal(t, `{"o":"https://example.com/test?param=abc","ts":1589327421000,"r":"https://www.google.com/search"}`, string(json))
}

func TestFromRequestReferrerSameHost(t *testing.T) {
	req := createRequest("https://example.com/test?param=abc", "https://example.com/search", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	event := FromRequest(req)
	fixTime(&event)

	assert.Equal(t, "https", event.Scheme)
	assert.Equal(t, "example.com", event.Host)
	assert.Equal(t, "", event.Referrer)
	assert.Equal(t, "Mozilla/5.0 (Windows NT 10.0; Win64; x64)", event.UserAgent)

	json, err := json.Marshal(event)
	assert.NoError(t, err)
	assert.Equal(t, `{"o":"https://example.com/test?param=abc","ts":1589327421000}`, string(json))
}

func TestFromRequestReferrerAbsent(t *testing.T) {
	req := createRequest("https://example.com/test?param=abc", "", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	event := FromRequest(req)
	fixTime(&event)

	assert.Equal(t, "https", event.Scheme)
	assert.Equal(t, "example.com", event.Host)
	assert.Equal(t, "", event.Referrer)
	assert.Equal(t, "Mozilla/5.0 (Windows NT 10.0; Win64; x64)", event.UserAgent)

	json, err := json.Marshal(event)
	assert.NoError(t, err)
	assert.Equal(t, `{"o":"https://example.com/test?param=abc","ts":1589327421000}`, string(json))
}

func TestSend(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		assert.NotEmpty(t, string(body))
		assert.Equal(t, Endpoint, r.URL.Path)
	}))
	defer svr.Close()

	req := createRequest(svr.URL+"/test?param=abc", "https://www.google.com/search", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	event := FromRequest(req)
	fixTime(&event)

	assert.Equal(t, svr.URL+Endpoint, event.Scheme+"://"+event.Host+Endpoint)

	err := Send(&event)
	assert.NoError(t, err)
}
