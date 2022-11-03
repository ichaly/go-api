package core

import (
	"encoding/json"
	"net/http"
)

var (
	OK    = newResult(http.StatusOK, "ok")                     // 通用成功
	ERROR = newResult(http.StatusInternalServerError, "error") // 通用错误
)

type result struct {
	Code    int         `json:"code"`    // 错误码
	Message string      `json:"message"` // 错误描述
	Data    interface{} `json:"data"`    // 返回数据
}

// WithMsg 自定义响应信息
func (res *result) WithMsg(message string) result {
	return result{
		Code:    res.Code,
		Message: message,
		Data:    res.Data,
	}
}

// WithData 追加响应数据
func (res *result) WithData(data interface{}) result {
	return result{
		Code:    res.Code,
		Message: res.Message,
		Data:    data,
	}
}

// ToString 返回 JSON 格式的错误详情
func (res *result) ToString() string {
	err := &struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}{
		Code:    res.Code,
		Message: res.Message,
		Data:    res.Data,
	}
	raw, _ := json.Marshal(err)
	return string(raw)
}

// 构造函数
func newResult(code int, msg string) *result {
	return &result{
		Code:    code,
		Message: msg,
		Data:    nil,
	}
}
