package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sechecker/pkg/sechecker"
	"time"
)

const eventFilePath = "event.json"

func main() {
	// 設定ファイルの読み込み

	//--- 前回のイベントファイルの読み込み
	var prevEventMetadata sechecker.MetaData
	// 比較をするためにゼロ値を nil ではなく空配列にする
	prevEventMetadata.Events = []sechecker.ScheduleEvent{}
	readEventFile(eventFilePath, &prevEventMetadata)

	//--- 現在発生しているイベントの取得
	var currentEventMetadata sechecker.MetaData
	body, err := getScheduleEvent(&currentEventMetadata)
	if err != nil {
		fmt.Printf("Get scheduleevent is error\n")
		os.Exit(1)
	}

	//--- 前回のイベントと現在発生しているイベントの比較
	eventState := diffEvent(prevEventMetadata, currentEventMetadata)

	if eventState == sechecker.StateNew {
		// アクションの実行
		doAction(currentEventMetadata)
	}

	//--- 現在のイベントをファイルに書き出し
	if err := writeEventFile(eventFilePath, currentEventMetadata); err != nil {
		fmt.Println("Write file error")
		os.Exit(1)
	}

	fmt.Printf("%s\n", body)

}

func readEventFile(filepath string, metadata *sechecker.MetaData) {
	// イベントファイルの読み込み
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		// ファイルが読み込めないときは何もしない
	} else {
		json.Unmarshal(data, &metadata)
	}

}

func writeEventFile(filepath string, metadata sechecker.MetaData) error {
	data, err := json.MarshalIndent(metadata, "", " ")
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

func getScheduleEvent(currentEventMetadata *sechecker.MetaData) ([]byte, error) {
	api, err := sechecker.NewScheduleEventAPI(nil, sechecker.MetadataServiceURI)
	if err != nil {
		fmt.Printf("Error initialize API")
		return nil, err
	}

	body, err := api.Run(currentEventMetadata)

	if err != nil {
		fmt.Printf("Error initialize API")
		return nil, err
	}
	return body, nil
}

func diffEvent(prevEventMetadata sechecker.MetaData, currentEventMetadata sechecker.MetaData) sechecker.EventState {
	var eventState sechecker.EventState
	// イベントの数で判断
	eventDiff := len(prevEventMetadata.Events) - len(currentEventMetadata.Events)
	switch {
	case len(prevEventMetadata.Events) == 0 && len(currentEventMetadata.Events) == 0:
		// ex) 0 -> 0
		eventState = sechecker.StateClosed
	case len(prevEventMetadata.Events) > 0 && len(currentEventMetadata.Events) == 0:
		// ex) 1 -> 0, 3 -> 0
	case eventDiff < 0:
		// ex) 0 -> 1, 1 -> 3
		eventState = sechecker.StateNew
	case eventDiff == 0:
		// ex) 1 -> 1, 3 -> 3
		eventState = sechecker.StateActive
	case eventDiff > 0:
		// ex) 3 -> 2, 2 -> 1
		eventState = sechecker.StateActive
	}
	return eventState
}

// アクションの実行
func doAction(currentEventMetadata sechecker.MetaData) error {

	// --- Pixelaへのポスト
	api, err := sechecker.NewPixelaClient("tsubasaxzzz", "scheduleevent", "organza-faun-weak")
	if err != nil {
		fmt.Println("Error")
	}

	err2 := api.PostEvent(currentEventMetadata)
	if err2 != nil {
		fmt.Println("Error")
	}

	// --- イベントファイルのコピー作成
	date := time.Now().Format("20060102150405")
	if err := writeEventFile(date+"_"+eventFilePath, currentEventMetadata); err != nil {
		fmt.Println("Write file error")
	}

	// --- Slack へのポスト

	// --- Log Analytics へのポスト
	return nil
}
