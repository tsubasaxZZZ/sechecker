package main

import (
	"fmt"
	"os"
	"sechecker/pkg/sechecker"
	"time"
)

const eventFilePath = "event.json"

func main() {
	// 設定ファイルの読み込み

	//--- 前回のイベントファイルの読み込み
	var prevEventMetadata sechecker.MetaData
	prevEventMetadata =	sechecker.ReadEventFile(eventFilePath)

	//--- 現在発生しているイベントの取得
	var currentEventMetadata sechecker.MetaData
	body, err := getScheduleEvent(&currentEventMetadata)
	if err != nil {
		fmt.Printf("Get scheduleevent is error\n")
		os.Exit(1)
	}

	//--- 前回のイベントと現在発生しているイベントの比較
	eventState := sechecker.DiffEvent(prevEventMetadata, currentEventMetadata)

	if eventState == sechecker.StateNew {
		// アクションの実行
		doAction(currentEventMetadata)
	}

	//--- 現在のイベントをファイルに書き出し
	if err := currentEventMetadata.WriteEventFile(eventFilePath); err != nil {
		fmt.Println("Write file error")
		os.Exit(1)
	}

	fmt.Printf("%s\n", body)

}

func getScheduleEvent(currentEventMetadata *sechecker.MetaData) ([]byte, error) {
	client, err := sechecker.New(nil, sechecker.MetadataServiceURI)
	if err != nil {
		fmt.Printf("Error initialize API")
		return nil, err
	}

	body, err := client.Run(currentEventMetadata)

	if err != nil {
		fmt.Printf("Error initialize API")
		return nil, err
	}
	return body, nil
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
	if err := currentEventMetadata.WriteEventFile(date+"_"+eventFilePath); err != nil {
		fmt.Println("Write file error")
	}

	// --- Slack へのポスト

	// --- Log Analytics へのポスト
	return nil
}
