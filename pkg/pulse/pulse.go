package pulse

import (
	"context"
	"net/http"

	"github.com/gobwas/ws"
)

type Pulse struct {
	opts *options
	hub  *hub
}

func New(opts ...Option) *Pulse {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	return &Pulse{
		opts: o,
		hub:  newHub(),
	}
}

func (p *Pulse) HandleRequest(w http.ResponseWriter, r *http.Request) error {

	// Origin 校验
	if p.opts.checkOrigin != nil {
		origin := r.Header.Get("Origin")
		if origin != "" && !p.opts.checkOrigin(origin) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return context.Canceled
		}
	}

	raw, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		return err
	}

	p.hub.allocate(context.Background(), p.opts, raw)

	return nil
}

// BroadcastBinary 广播二进制帧
// 注意这里默认每个连接都会 copy 一份（Write 内部 copy）
// 若你要极致优化：可以增加 BroadcastNoCopy + 强约束 msg 不可复用/修改
func (p *Pulse) BroadcastBinary(msg []byte) {
	p.hub.broadcastBinary(msg)
}

// BroadcastText 广播文本帧
func (p *Pulse) BroadcastText(msg []byte) {
	p.hub.broadcastText(msg)
}
