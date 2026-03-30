package cluster

import (
	"github.com/byteweap/meta/component/broker"
	"github.com/byteweap/meta/pkg/conv"
)

const (
	FieldName_Uid         = "uid"
	FieldName_Event       = "event"
	FieldName_Reply       = "reply"
	FieldName_FromService = "from_service"
	FieldName_ToService   = "to_service"
)

// BuildHeader 构建必备请求头
func BuildHeader(uid int64, event Event, reply, fromService, toService string) broker.Header {
	return broker.Header{
		FieldName_Uid:         []string{conv.String(uid)},
		FieldName_Event:       []string{string(event)},
		FieldName_Reply:       []string{reply},
		FieldName_FromService: []string{fromService},
		FieldName_ToService:   []string{toService},
	}
}

// GetUidBy 从请求头中获取用户ID
func GetUidBy(header broker.Header) int64 {
	return conv.Int64(header.Get(FieldName_Uid))
}

// GetEventBy 从请求头中获取事件类型
func GetEventBy(header broker.Header) Event {
	return Event(header.Get(FieldName_Event))
}

// GetReplyBy 从请求头中获取回复信息
func GetReplyBy(header broker.Header) string {
	return header.Get(FieldName_Reply)
}

// GetFromServiceBy 从请求头中获取来源服务
func GetFromServiceBy(header broker.Header) string {
	return header.Get(FieldName_FromService)
}

// GetToServiceBy 从请求头中获取目标服务
func GetToServiceBy(header broker.Header) string {
	return header.Get(FieldName_ToService)
}
