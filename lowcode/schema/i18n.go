package schema

import "encoding/json"

type Lang string
type LangDict map[string]string

type I18NSchema struct {
	DataRaw json.RawMessage   `json:"i18n"`
	Data    map[Lang]LangDict `json:"-"`
}

type I18NData struct {
	TypeName string `json:"type"`
	Key      string `json:"key"`
}

func (in *I18NData) NodeDataType() string {
	return "I18NData"
}
