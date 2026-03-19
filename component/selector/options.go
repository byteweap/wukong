package selector

// SelectOptions 是选择参数
type SelectOptions struct {
	NodeFilters []NodeFilter
}

// SelectOption 选择器选项
type SelectOption func(*SelectOptions)

// WithNodeFilter 设置过滤器
func WithNodeFilter(fn ...NodeFilter) SelectOption {
	return func(opts *SelectOptions) {
		opts.NodeFilters = fn
	}
}
