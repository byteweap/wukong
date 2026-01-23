package zookeeper

import (
	"encoding/json"

	"github.com/byteweap/wukong/component/registry"
)

// marshal 将服务实例编码为 JSON
func marshal(si *registry.ServiceInstance) ([]byte, error) {
	return json.Marshal(si)
}

// unmarshal 将 JSON 解码为服务实例
func unmarshal(data []byte) (si *registry.ServiceInstance, err error) {
	err = json.Unmarshal(data, &si)
	return
}
