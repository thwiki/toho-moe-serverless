package handler

import (
	"net/http"
	"net/url"
	"strings"

	utils "github.com/thwiki/toho-moe-serverless/utils"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	slug := strings.Trim(r.URL.Path, "/ ")
	userId := utils.GetUserIP(r)
	userAgent := r.Header.Get("User-Agent")

	header := w.Header()
	header.Set("Content-Type", "application/json; charset=utf-8")
	header.Set("Cache-Control", "public, max-age=0, s-maxage=0, must-revalidate")

	date := utils.GetDate()
	header.Set("Last-Modified", date.Format(http.TimeFormat))

	shortUrls := utils.GetShortUrls()
	shortUrl, ok := shortUrls[slug]

	var value float32
	if ok {
		value = 2.0
	} else {
		value = 1.0
	}

	go utils.SendAnalytics(userId, userAgent, value, slug)

	if !ok {
		http.NotFound(w, r)
		return
	}

	target, err := url.ParseRequestURI(shortUrl.Url)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	targetQuery := target.Query()

	for key, value := range query {
		targetQuery.Del(key)

		for i := 0; i < len(value); i++ {
			targetQuery.Add(key, value[i])
		}
	}

	target.RawQuery = targetQuery.Encode()

	http.Redirect(w, r, target.String(), http.StatusFound)
}
