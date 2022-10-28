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
	TransactionId      string              `json:"transaction_id"`
	Currency           string              `json:"currency"`
	Value              float32             `json:"value"`
	Items              []AnalyticEventItem `json:"items"`
	NonPersonalizedAds bool                `json:"non_personalized_ads"`
	EngagementTimeMsec string              `json:"engagement_time_msec"`
}

type AnalyticEventItem struct {
	ItemId   string  `json:"item_id"`
	Currency string  `json:"currency"`
	Value    float32 `json:"value"`
	Quantity float32 `json:"quantity"`
}

var (
	endpoint      = os.Getenv("GA_ENDPOINT")
	apiSecret     = os.Getenv("GA_API_SECRET")
	measurementId = os.Getenv("GA_MEASUREMENT_ID")
)

func SendAnalytics(userId string, userAgent string, value float32, slug string) {
	endpoint := endpoint + "?api_secret=" + apiSecret + "&measurement_id=" + measurementId

	id := uuid.New().String()

	analyticData := AnalyticData{}
	analyticData.ClientId = id
	analyticData.UserId = userId
	analyticData.Events = make([]AnalyticEvent, 1)

	analyticEvent := AnalyticEvent{
		Name: "purchase",
		Params: AnalyticEventParams{
			TransactionId:      id,
			Currency:           "CNY",
			Value:              value,
			Items:              make([]AnalyticEventItem, 1),
			NonPersonalizedAds: true,
			EngagementTimeMsec: "1",
		},
	}

	analyticItem := AnalyticEventItem{
		ItemId:   slug,
		Currency: "CNY",
		Value:    value,
		Quantity: 1,
	}

	analyticData.Events[0] = analyticEvent
	analyticEvent.Params.Items[0] = analyticItem

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
