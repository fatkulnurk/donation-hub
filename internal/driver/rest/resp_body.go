package rest

import "time"

type RespBodyOk struct {
	OK        bool        `json:"ok"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"ts"`
}

func NewResponseOk(data interface{}) *RespBodyOk {
	return &RespBodyOk{
		OK:        true,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}
