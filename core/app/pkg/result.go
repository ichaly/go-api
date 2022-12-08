package core

import (
	"encoding/json"
	"net/http"
)

var (
	OK    = NewResult(http.StatusOK)                  // 通用成功
	ERROR = NewResult(http.StatusInternalServerError) // 通用错误
)

type result struct {
	Code   int                      `json:"code"`             // 错误码
	Data   interface{}              `json:"data,omitempty"`   // 返回数据
	Errors []map[string]interface{} `json:"errors,omitempty"` // 错误信息
}

// WithError 自定义错误信息
func (res *result) WithError(errors ...error) result {
	if errors != nil && len(errors) > 0 {
		for _, e := range errors {
			res.Errors = append(res.Errors, map[string]interface{}{"message": e.Error()})
		}
	}
	return result{
		Code:   res.Code,
		Errors: res.Errors,
	}
}

// WithData 追加响应数据
func (res *result) WithData(data interface{}) result {
	return result{
		Code: res.Code,
		Data: data,
	}
}

// ToString 返回 JSON 格式的错误详情
func (res *result) ToString() string {
	err := &struct {
		Code   int                      `json:"code"`
		Data   interface{}              `json:"data,omitempty"`
		Errors []map[string]interface{} `json:"errors,omitempty"`
	}{
		Code:   res.Code,
		Data:   res.Data,
		Errors: res.Errors,
	}
	raw, _ := json.Marshal(err)
	return string(raw)
}

// NewResult 构造函数
func NewResult(code int) *result {
	return &result{
		Code: code,
		Data: nil,
	}
}
