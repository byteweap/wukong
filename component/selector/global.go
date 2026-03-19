package selector

// globalSelector 保存全局选择器
var globalSelector = &wrapSelector{}

var _ Builder = (*wrapSelector)(nil)

// wrapSelector 包装 Selector 用于覆盖全局实现
type wrapSelector struct{ Builder }

// GlobalSelector 返回全局 Selector 构建器
func GlobalSelector() Builder {
	if globalSelector.Builder != nil {
		return globalSelector
	}
	return nil
}

// SetGlobalSelector 设置全局 Selector 构建器
func SetGlobalSelector(builder Builder) {
	globalSelector.Builder = builder
}
