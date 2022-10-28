package handler

import (
	"net/http"
	"net/url"
	"os"

	utils "github.com/thwiki/toho-moe-serverless/utils"
)

var (
	thbcdUrl = os.Getenv("THBCD_URL")
)

func THBCD(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	number := query.Get("number")

	userId := utils.GetUserIP(r)
	userAgent := r.Header.Get("User-Agent")

	header := w.Header()
	header.Set("Content-Type", "application/json; charset=utf-8")
	header.Set("Cache-Control", "public, max-age=0, s-maxage=0, must-revalidate")

	date := utils.GetDate()
	header.Set("Last-Modified", date.Format(http.TimeFormat))

	go utils.SendAnalytics(userId, userAgent, 3.0, name+"#"+number)

	target, err := url.ParseRequestURI(thbcdUrl)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	targetQuery := target.Query()
	targetQuery.Add("e", name+"#"+number)

	target.RawQuery = targetQuery.Encode()

	http.Redirect(w, r, target.String(), http.StatusFound)
}
