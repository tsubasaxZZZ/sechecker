package sechecker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Actioner interface {
	PostEvent(metadata MetaData) error
}

type ActionConfig struct {
	Name   string   `json:"Name"`
	Type   string   `json:"Type"`
	Config Actioner `json:"Config"`
}

type Configs struct {
	ActionConfigs []ActionConfig `json:"ActionConfig"`
}

func (c *Configs) UnmarshalJSON(b []byte) error {
	type alias Configs
	a := struct {
		*alias
	}{
		alias: (*alias)(c),
	}
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	return nil
}

var (
	_configObj = map[string]Actioner{
		// Pixela 用のコンフィグ
		"Pixela": &PixelaConfig{},
	}
)

func (c *ActionConfig) UnmarshalJSON(b []byte) error {
	type alias ActionConfig
	a := struct {
		Config json.RawMessage `json:"Config"`
		*alias
	}{
		alias: (*alias)(c),
	}
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}

	// Type から型情報を判別しマップから取得
	// 取得できないときは nil
	o, ok := _configObj[c.Type]
	if ok {
		if err := json.Unmarshal(a.Config, &o); err != nil {
			return err
		}
	}
	c.Config = o
	return nil
}

func ReadConfig(filepath string, c *Configs) error {
	// ファイル存在チェック
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		// ファイルが存在しない場合は空コンフィグを生成
		c := Configs{}
		c.ActionConfigs = []ActionConfig{}
		WriteConfig(filepath, c)
	}

	// イベントファイルの読み込み
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("ReadConfig io error: %s", err)
	}

	if err := json.Unmarshal(data, c); err != nil {
		return fmt.Errorf("ReadConfig Unmarshal error: %s", err)
	}

	return nil
}

func WriteConfig(filepath string, c Configs) error {
	data, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath, data, 0644); err != nil {
		return err
	}
	return nil
}
