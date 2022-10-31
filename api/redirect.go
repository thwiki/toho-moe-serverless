package handler

import (
	"net/http"
	"net/url"
	"strings"

	utils "github.com/thwiki/toho-moe-serverless/utils"
	"github.com/thwiki/toho-moe-serverless/vago"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	slug := strings.Trim(r.URL.Path, "/ ")

	header := w.Header()
	header.Set("Content-Type", "application/json; charset=utf-8")
	header.Set("Cache-Control", "public, max-age=0, s-maxage=0, must-revalidate")

	date := utils.GetDate()
	header.Set("Last-Modified", date.Format(http.TimeFormat))

	shortUrls := utils.GetShortUrls()
	shortUrl, ok := shortUrls[slug]

	event := vago.FromRequest(r)

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

	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}

	vago.Send(&event)
}
