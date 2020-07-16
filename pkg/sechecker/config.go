package sechecker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Commands []string `json:"command"`
}

func ReadConfig(filepath string, c *Config) error {
	// ファイル存在チェック
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		// ファイルが存在しない場合は空コンフィグを生成
		c := Config{}
		WriteConfig(filepath, c)
	}

	// コンフィグファイルの読み込み
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("ReadConfig io error: %s", err)
	}

	if err := json.Unmarshal(data, c); err != nil {
		return fmt.Errorf("ReadConfig Unmarshal error: %s", err)
	}

	return nil
}

func WriteConfig(filepath string, c Config) error {
	data, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath, data, 0644); err != nil {
		return err
	}
	return nil
}
