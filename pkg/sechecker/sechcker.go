package sechecker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const MetadataServiceURI = "http://169.254.169.254/metadata/scheduledevents?api-version=2017-11-01"

type EventState int

const (
	// イベントが終了した(予定されているイベントがない)状態を表す
	StateClosed EventState = iota
	// 新しいイベントが発生した状態を表す
	StateNew
	// イベントが継続中の状態を表す
	StateActive
)

type ScheduleEventAPI struct {
	Client  *http.Client
	BaseURL string
}

type ScheduleEvent struct {
	EventID      string   `json:"EventId"`
	EventType    string   `json:"EventType"`
	ResourceType string   `json:"ResourceType"`
	Resources    []string `json:"Resources"`
	EventStatus  string   `json:"EventStatus"`
	NotBefore    string   `json:"NotBefore"`
}

type MetaData struct {
	EventState          EventState
	DocumentIncarnation int             `json:"DocumentIncarnation"`
	Events              []ScheduleEvent `json:"Events"`
}

func NewScheduleEventAPI(client *http.Client, baseURL string) (*ScheduleEventAPI, error) {
	hc := client
	if hc == nil {
		hc = &http.Client{Timeout: time.Duration(120) * time.Second}
	}

	if baseURL == "" {
		baseURL = MetadataServiceURI
	}
	return &ScheduleEventAPI{
		Client:  hc,
		BaseURL: baseURL,
	}, nil
}

func (api *ScheduleEventAPI) Run(metadata *MetaData) ([]byte, error) {
	req, err1 := http.NewRequest("GET", api.BaseURL, nil)
	if err1 != nil {
		fmt.Println("http.NewRequest error")
		return nil, err1
	}
	req.Header.Set("Metadata", "true")

	resp, err2 := api.Client.Do(req)
	if err2 != nil {
		fmt.Println("client.Do Error")
		return nil, err2
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	jsonBytes := ([]byte)(byteArray)

	err3 := json.Unmarshal(jsonBytes, metadata)
	if err3 != nil {
		fmt.Println("json.Unmarshal Error")
		return nil, err3
	}
	return jsonBytes, nil
}
