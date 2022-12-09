package json

import (
	jsoniter "github.com/json-iterator/go"
)

var (
	json                = jsoniter.ConfigCompatibleWithStandardLibrary
	Marshal             = json.Marshal
	Unmarshal           = json.Unmarshal
	MarshalToString     = json.MarshalToString
	UnmarshalFromString = json.UnmarshalFromString
)
