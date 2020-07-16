package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sechecker/pkg/sechecker"
	"strings"
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
	var config sechecker.Config
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
	log.Printf("EventState=%d\n", eventState)

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
func doAction(currentEventMetadata sechecker.MetaData, c sechecker.Config) (int, error) {

	// 実行されたアクションの数
	var acctionCount int

	for _, command := range c.Commands {
		log.Println(command)
		cmd := exec.Command("sh", "-c", command)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			log.Printf("%s:%s", strings.TrimRight(strings.TrimRight(stderr.String(), "\n"), "\n"), err)
			continue
		}
		log.Println(stdout.String())
		acctionCount++
	}

	// --- イベントファイルのコピー作成
	date := time.Now().Format("20060102150405")
	if err := currentEventMetadata.WriteEventFile(date + "_" + eventFilePath); err != nil {
		fmt.Println("Write file error")
	}

	return acctionCount, nil
}
