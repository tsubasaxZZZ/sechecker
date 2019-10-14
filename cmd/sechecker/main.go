package main

import (
	"fmt"
	"os"
	"sechecker/pkg/sechecker"
	"time"
)

const eventFilePath = "event.json"
const configFilePath = "config.json"

// exit codes
const (
	ExitCodeOK  = 0
	ExitCodeErr = 3
)

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(ExitCodeErr)
		return
	}
	os.Exit(ExitCodeOK)
}

func run() error {

	// 設定ファイルの読み込み
	var config sechecker.Configs
	if err := sechecker.ReadConfig(configFilePath, &config); err != nil {
		return err
	}

	//--- 前回のイベントファイルの読み込み
	var prevEventMetadata sechecker.MetaData
	prevEventMetadata = sechecker.ReadEventFile(eventFilePath)

	//--- 現在発生しているイベントの取得
	var currentEventMetadata sechecker.MetaData
	body, err := getScheduleEvent(&currentEventMetadata)
	if err != nil {
		return err
	}

	//--- 前回のイベントと現在発生しているイベントの比較
	eventState := sechecker.DiffEvent(prevEventMetadata, currentEventMetadata)

	if eventState == sechecker.StateNew {
		// アクションの実行
		doAction(currentEventMetadata, config)
	}

	//--- 現在のイベントをファイルに書き出し
	if err := currentEventMetadata.WriteEventFile(eventFilePath); err != nil {
		return err
	}

	fmt.Printf("%s\n", body)
	return nil

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
func doAction(currentEventMetadata sechecker.MetaData, c sechecker.Configs) error {

	for _, actionConfig := range c.ActionConfigs {
		switch v := actionConfig.Config.(type) {
		case *sechecker.PixelaConfig:
			api := sechecker.NewPixelaClient(v.UserID, v.GraphID, v.Secret)
			err := api.PostEvent(currentEventMetadata)
			if err != nil {
				// エラーになっても次のアクションを実行
				fmt.Printf("%s", err)
			}
		}

	}

	// --- イベントファイルのコピー作成
	date := time.Now().Format("20060102150405")
	if err := currentEventMetadata.WriteEventFile(date + "_" + eventFilePath); err != nil {
		fmt.Println("Write file error")
	}

	// --- Slack へのポスト

	// --- Log Analytics へのポスト
	return nil
}
