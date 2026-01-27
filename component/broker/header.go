package broker

// Header 是消息头.实现可以选择性支持
type Header map[string][]string

// Add 添加键值对到 header，追加到已有值后面 (区分大小写)
func (h Header) Add(key, value string) {
	h[key] = append(h[key], value)
}

// Set 设置 header 键的值，替换已有值 (区分大小写)
func (h Header) Set(key, value string) {
	h[key] = []string{value}
}

// Get 获取指定键的第一个值 (区分大小写)
func (h Header) Get(key string) string {
	if h == nil {
		return ""
	}
	if v := h[key]; v != nil {
		return v[0]
	}
	return ""
}

// Values 获取指定键的所有值 (区分大小写)
func (h Header) Values(key string) []string {
	return h[key]
}

// Del 删除指定键 (区分大小写)
func (h Header) Del(key string) {
	delete(h, key)
}
