package sechecker

import "encoding/json"

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
