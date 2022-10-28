package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/google/uuid"
)

type AnalyticData struct {
	ClientId string          `json:"client_id"`
	UserId   string          `json:"user_id"`
	Events   []AnalyticEvent `json:"events"`
}

type AnalyticEvent struct {
	Name   string              `json:"name"`
	Params AnalyticEventParams `json:"params"`
}

type AnalyticEventParams struct {
	ContentType        string `json:"content_type"`
	ItemId             string `json:"item_id"`
	NonPersonalizedAds bool   `json:"non_personalized_ads"`
	EngagementTimeMsec string `json:"engagement_time_msec"`
}

var (
	endpoint      = os.Getenv("GA_ENDPOINT")
	apiSecret     = os.Getenv("GA_API_SECRET")
	measurementId = os.Getenv("GA_MEASUREMENT_ID")
)

func SendAnalytics(userId string, userAgent string, actionType string, slug string) {
	endpoint := endpoint + "?api_secret=" + apiSecret + "&measurement_id=" + measurementId

	analyticData := AnalyticData{}
	analyticData.ClientId = uuid.New().String()
	analyticData.UserId = userId

	analyticEvent := AnalyticEvent{
		Name: "select_content",
		Params: AnalyticEventParams{
			ContentType:        actionType,
			ItemId:             slug,
			NonPersonalizedAds: true,
			EngagementTimeMsec: "1",
		},
	}

	analyticData.Events = append(analyticData.Events, analyticEvent)

	json, err := json.Marshal(analyticData)

	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(json))

	if err != nil {
		return
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return
	}

	defer resp.Body.Close()
	io.Copy(ioutil.Discard, resp.Body)
}
