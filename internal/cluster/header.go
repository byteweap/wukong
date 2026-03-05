package cluster

import (
	"github.com/byteweap/wukong/component/broker"
	"github.com/byteweap/wukong/pkg/conv"
)

const (
	FieldName_Uid   = "uid"
	FieldName_Event = "event"
	FieldName_Reply = "reply"
)

// BuildHeader 构建必备请求头
func BuildHeader(uid int64, event Event, reply string) broker.Header {
	return broker.Header{
		FieldName_Uid:   []string{conv.String(uid)},
		FieldName_Event: []string{string(event)},
		FieldName_Reply: []string{reply},
	}
}

// GetUidByHeader 从请求头中获取用户ID
func GetUidByHeader(header broker.Header) int64 {
	return conv.Int64(header.Get(FieldName_Uid))
}

// GetEventByHeader 从请求头中获取事件类型
func GetEventByHeader(header broker.Header) Event {
	return Event(header.Get(FieldName_Event))
}

// GetReplyByHeader 从请求头中获取回复信息
func GetReplyByHeader(header broker.Header) string {
	return header.Get(FieldName_Reply)
}
