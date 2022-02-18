package json

import (
	"github.com/draymonders/bingo/log"
	jsoniter "github.com/json-iterator/go"
)

var (
	api  jsoniter.API
)

func init() {
	json := jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		UseNumber:              true,
	}
	api = json.Froze()
}

func ObjToJson(obj interface{}) string {
	s, err := api.MarshalToString(obj)
	if err != nil {
		log.Error("marshal err. %+v", err)
	}
	return s
}

func JsonToObj(str string, obj interface{}) error {
	return api.UnmarshalFromString(str, obj)
}
