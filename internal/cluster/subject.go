package cluster

// Subject 组装主题
// subject = 前缀.发送方AppName.接收方AppName.接收方节点ID
func Subject(prefix, fromApp, toApp, toAppID string) string {
	if prefix == "" {
		return fromApp + "." + toApp + "." + toAppID
	}
	return prefix + "." + fromApp + "." + toApp + "." + toAppID
}
