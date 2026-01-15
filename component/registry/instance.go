package registry

import (
	"fmt"
	"sort"
)

// ServiceInstance 服务发现系统中的服务实例
type ServiceInstance struct {
	// ID 注册时的唯一实例ID
	ID string `json:"id"`
	// Name 注册时的服务名
	Name string `json:"name"`
	// Version 编译版本
	Version string `json:"version"`
	// Metadata 服务实例关联的键值对元数据
	Metadata map[string]string `json:"metadata"`
	// Endpoints 服务实例的端点地址
	// 格式示例:
	//   http://127.0.0.1:8000?isSecure=false
	//   grpc://127.0.0.1:9000?isSecure=false
	Endpoints []string `json:"endpoints"`
}

// String 返回服务实例的字符串表示
func (i *ServiceInstance) String() string {
	return fmt.Sprintf("%s-%s", i.Name, i.ID)
}

// Equal 判断两个服务实例是否相等
func (i *ServiceInstance) Equal(o any) bool {
	if i == nil && o == nil {
		return true
	}

	if i == nil || o == nil {
		return false
	}

	t, ok := o.(*ServiceInstance)
	if !ok {
		return false
	}

	if len(i.Endpoints) != len(t.Endpoints) {
		return false
	}

	sort.Strings(i.Endpoints)
	sort.Strings(t.Endpoints)
	for j := 0; j < len(i.Endpoints); j++ {
		if i.Endpoints[j] != t.Endpoints[j] {
			return false
		}
	}

	if len(i.Metadata) != len(t.Metadata) {
		return false
	}

	for k, v := range i.Metadata {
		if v != t.Metadata[k] {
			return false
		}
	}

	return i.ID == t.ID && i.Name == t.Name && i.Version == t.Version
}
