package core

import "encoding/json"

type Variable map[string]interface{}

func (my Variable) Marshal() json.RawMessage {
	if res, err := json.Marshal(my); err == nil {
		return res
	}
	return nil
}
