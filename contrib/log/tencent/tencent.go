package tencent

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	cls "github.com/tencentcloud/tencentcloud-cls-sdk-go"
	"google.golang.org/protobuf/proto"

	"github.com/byteweap/wukong/component/log"
)

type Logger interface {
	log.Logger

	GetProducer() *cls.AsyncProducerClient
	Close() error
}

type tencentLog struct {
	producer *cls.AsyncProducerClient
	opts     *options
}

var _ Logger = (*tencentLog)(nil)

func NewLogger(options ...Option) (Logger, error) {
	opts := defaultOptions()
	for _, o := range options {
		o(opts)
	}
	producerConfig := cls.GetDefaultAsyncProducerClientConfig()
	producerConfig.AccessKeyID = opts.accessKey
	producerConfig.AccessKeySecret = opts.accessSecret
	producerConfig.Endpoint = opts.endpoint
	producerInst, err := cls.NewAsyncProducerClient(producerConfig)
	if err != nil {
		return nil, err
	}
	producerInst.Start()
	return &tencentLog{
		producer: producerInst,
		opts:     opts,
	}, nil
}

func (log *tencentLog) GetProducer() *cls.AsyncProducerClient {
	return log.producer
}

func (log *tencentLog) Close() error {
	return log.producer.Close(5000)
}

func (log *tencentLog) Log(level log.Level, kvs ...any) error {
	contents := make([]*cls.Log_Content, 0, len(kvs)/2+1)

	contents = append(contents, &cls.Log_Content{
		Key:   newString(level.Key()),
		Value: newString(level.String()),
	})
	for i := 0; i < len(kvs); i += 2 {
		contents = append(contents, &cls.Log_Content{
			Key:   newString(toString(kvs[i])),
			Value: newString(toString(kvs[i+1])),
		})
	}

	logInst := &cls.Log{
		Time:     proto.Int64(time.Now().Unix()),
		Contents: contents,
	}
	return log.producer.SendLog(log.opts.topicID, logInst, nil)
}

func newString(s string) *string {
	return &s
}

// toString convert any type to string
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
