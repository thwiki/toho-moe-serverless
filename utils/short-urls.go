package utils

import (
	"encoding/json"
	"time"
)

type ShortUrl struct {
	Slug string `json:"slug"`
	Url  string `json:"url"`
}

type ShortUrls map[string]ShortUrl

func UnmarshalShortUrls(data []byte) (m ShortUrls) {
	if err := json.Unmarshal(data, &m); err != nil {
		return
	}
	return
}

func GetShortUrls() ShortUrls {
	return UnmarshalShortUrls(data)
}

func GetDate() time.Time {
	return time.UnixMilli(int64(date))
}
