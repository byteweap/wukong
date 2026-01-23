package etcd

import (
	"encoding/json"

	"github.com/byteweap/wukong/component/registry"
)

// marshal 将服务实例编码为 JSON 字符串
func marshal(si *registry.ServiceInstance) (string, error) {
	data, err := json.Marshal(si)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// unmarshal 将 JSON 解码为服务实例
func unmarshal(data []byte) (si *registry.ServiceInstance, err error) {
	err = json.Unmarshal(data, &si)
	return
}
