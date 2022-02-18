package json

import (
	"github.com/draymonders/bingo/log"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func ObjToJson(obj interface{}) string {
	s, err := json.MarshalToString(obj)
	if err != nil {
		log.Error("marshal err. %+v", err)
	}
	return s
}

func JsonToObj(str string, obj interface{}) error {
	return json.UnmarshalFromString(str, obj)
}
