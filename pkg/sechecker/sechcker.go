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
	// StateClosed はイベントが終了した(予定されているイベントがない)状態を表す
	StateClosed EventState = iota
	// StateNew は新しいイベントが発生した状態を表す
	StateNew
	// StateActive はイベントが継続中の状態を表す
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

func New(client *http.Client, baseURL string) (*ScheduleEventAPI, error) {
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

func ReadEventFile(filepath string) MetaData {
	var metadata MetaData
	// 比較をするためにゼロ値を nil ではなく空配列にする
	metadata.Events = []ScheduleEvent{}

	// イベントファイルの読み込み
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		// ファイルが読み込めないときは何もしない
	} else {
		json.Unmarshal(data, &metadata)
	}
	return metadata

}
func (m *MetaData) WriteEventFile(filepath string) error {
	data, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		fmt.Println(err)
		return err
	}
	if err := ioutil.WriteFile(filepath, data, 0644); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func DiffEvent(prevEventMetadata MetaData, currentEventMetadata MetaData) EventState {
	var eventState EventState
	// イベントの数で判断
	eventDiff := len(prevEventMetadata.Events) - len(currentEventMetadata.Events)
	switch {
	case len(prevEventMetadata.Events) == 0 && len(currentEventMetadata.Events) == 0:
		// ex) 0 -> 0
		eventState = StateClosed
	case len(prevEventMetadata.Events) > 0 && len(currentEventMetadata.Events) == 0:
		// ex) 1 -> 0, 3 -> 0
		eventState = StateClosed
	case eventDiff < 0:
		// ex) 0 -> 1, 1 -> 3
		eventState = StateNew
	case eventDiff == 0:
		// ex) 1 -> 1, 3 -> 3
		eventState = StateActive
	case eventDiff > 0:
		// ex) 3 -> 2, 2 -> 1
		eventState = StateActive
	}
	return eventState
}
