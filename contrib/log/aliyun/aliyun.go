package aliyun

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"google.golang.org/protobuf/proto"

	"github.com/byteweap/wukong/component/log"
)

// Logger 阿里云日志记录器接口
type Logger interface {
	log.Logger

	GetProducer() *producer.Producer
	Close() error
}

type aliyunLog struct {
	producer *producer.Producer
	opts     *options
}

var _ Logger = (*aliyunLog)(nil)

// NewAliyunLog 创建阿里云日志记录器
func NewAliyunLog(options ...Option) (Logger, error) {
	opts := defaultOptions()
	for _, o := range options {
		o(opts)
	}

	producerConfig := producer.GetDefaultProducerConfig()
	producerConfig.Endpoint = opts.endpoint
	producerConfig.AccessKeyID = opts.accessKey
	producerConfig.AccessKeySecret = opts.accessSecret
	producerInst, err := producer.NewProducer(producerConfig)
	if err != nil {
		return nil, err
	}
	producerInst.Start()

	return &aliyunLog{
		opts:     opts,
		producer: producerInst,
	}, nil
}

// GetProducer 获取生产者
func (a *aliyunLog) GetProducer() *producer.Producer {
	return a.producer
}

// Close 关闭生产者
func (a *aliyunLog) Close() error {
	return a.producer.Close(5000)
}

// Log 发送日志
func (a *aliyunLog) Log(level log.Level, kvs ...any) error {
	contents := make([]*sls.LogContent, 0, len(kvs)/2+1)

	contents = append(contents, &sls.LogContent{
		Key:   newString(level.Key()),
		Value: newString(level.String()),
	})
	for i := 0; i < len(kvs); i += 2 {
		contents = append(contents, &sls.LogContent{
			Key:   newString(toString(kvs[i])),
			Value: newString(toString(kvs[i+1])),
		})
	}

	logInst := &sls.Log{
		Time:     proto.Uint32(uint32(time.Now().Unix())),
		Contents: contents,
	}
	return a.producer.SendLog(a.opts.project, a.opts.logstore, "", "", logInst)
}

// newString 转换为字符串指针
func newString(s string) *string {
	return &s
}

// toString 转换为字符串
func toString(v any) string {
	var key string
	if v == nil {
		return key
	}
	switch v := v.(type) {
	case float64:
		key = strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		key = strconv.FormatFloat(float64(v), 'f', -1, 32)
	case int:
		key = strconv.Itoa(v)
	case uint:
		key = strconv.FormatUint(uint64(v), 10)
	case int8:
		key = strconv.Itoa(int(v))
	case uint8:
		key = strconv.FormatUint(uint64(v), 10)
	case int16:
		key = strconv.Itoa(int(v))
	case uint16:
		key = strconv.FormatUint(uint64(v), 10)
	case int32:
		key = strconv.Itoa(int(v))
	case uint32:
		key = strconv.FormatUint(uint64(v), 10)
	case int64:
		key = strconv.FormatInt(v, 10)
	case uint64:
		key = strconv.FormatUint(v, 10)
	case string:
		key = v
	case bool:
		key = strconv.FormatBool(v)
	case []byte:
		key = string(v)
	case fmt.Stringer:
		key = v.String()
	default:
		newValue, _ := json.Marshal(v)
		key = string(newValue)
	}
	return key
}
