package vago

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type VAEvent struct {
	Url       string `json:"o"`
	Timestamp int64  `json:"ts"`
	Referrer  string `json:"r,omitempty"`
	IP        string `json:"-"`
	Scheme    string `json:"-"`
	Host      string `json:"-"`
	UserAgent string `json:"-"`
}

var (
	VA_API  = os.Getenv("VA_API")
	LOG_API = os.Getenv("LOG_API")
)

func FromRequest(r *http.Request) VAEvent {
	url := r.URL

	ip := getUserIP(r)

	host := url.Host

	if host == "" {
		host = r.Header.Get("Host")
	}
	if host == "" {
		host = os.Getenv("VERCEL_URL")
	}

	scheme := url.Scheme

	if scheme == "" {
		scheme = "https"
	}

	userAgent := r.Header.Get("User-Agent")
	referrer := r.Header.Get("Referer")
	timestamp := time.Now().UnixMilli()

	url.Host = host
	url.Scheme = scheme

	event := VAEvent{
		Url:       url.String(),
		Timestamp: timestamp,
		IP:        ip,
		Scheme:    scheme,
		Host:      host,
		UserAgent: userAgent,
	}

	if !strings.Contains(referrer, host) {
		event.Referrer = referrer
	}

	return event
}

func Send(event *VAEvent) (err error) {
	json, err := json.Marshal(event)

	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", VA_API, bytes.NewBuffer(json))

	if err != nil {
		return
	}

	req.Header.Set("X-Real-IP", event.IP)
	req.Header.Set("X-Forwarded-For", event.IP)
	req.Header.Set("User-Agent", event.UserAgent)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return
	}

	defer resp.Body.Close()
	log, err := io.ReadAll(resp.Body)

	if err != nil {
		return
	}

	req2, err2 := http.NewRequest("POST", LOG_API, bytes.NewBuffer(log))

	if err2 != nil {
		return
	}

	req2.Header.Set("Content-Type", "application/json")

	client2 := &http.Client{}
	resp2, err := client2.Do(req2)

	io.Copy(io.Discard, resp2.Body)

	return
}

func getUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-IP")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
