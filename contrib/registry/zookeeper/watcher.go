package zookeeper

import (
	"context"
	"errors"
	"path"
	"sync/atomic"

	"github.com/byteweap/wukong/component/registry"
	"github.com/go-zookeeper/zk"
)

var ErrWatcherStopped = errors.New("watcher stopped")

type watcher struct {
	ctx    context.Context
	event  chan zk.Event
	conn   *zk.Conn
	cancel context.CancelFunc

	first atomic.Bool
	// ZooKeeper 路径前缀，用于过滤或识别被监听节点
	prefix string
	// ZooKeeper 中被监听的服务名
	serviceName string
}

var _ registry.Watcher = (*watcher)(nil)

func newWatcher(ctx context.Context, prefix, serviceName string, conn *zk.Conn) (*watcher, error) {
	w := &watcher{conn: conn, event: make(chan zk.Event, 1), prefix: prefix, serviceName: serviceName}
	w.ctx, w.cancel = context.WithCancel(ctx)
	go w.watch(w.ctx)
	return w, nil
}

func (w *watcher) watch(ctx context.Context) {
	for {
		// 单次 watch 只对一次事件有效，需要循环持续监听
		_, _, ch, err := w.conn.ChildrenW(w.prefix)
		if err != nil {
			// 目标服务节点尚未创建
			if errors.Is(err, zk.ErrNoNode) {
				// 监听节点是否存在
				_, _, ch, err = w.conn.ExistsW(w.prefix)
			}
			if err != nil {
				w.event <- zk.Event{Err: err}
				continue
			}
		}
		select {
		case <-ctx.Done():
			return
		case ev := <-ch:
			w.event <- ev
		}
	}
}

func (w *watcher) Next() ([]*registry.ServiceInstance, error) {
	// TODO: 多次调用 Next 可能导致实例信息不一致
	if w.first.CompareAndSwap(false, true) {
		return w.getServices()
	}
	select {
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	case e := <-w.event:
		if e.State == zk.StateDisconnected {
			return nil, ErrWatcherStopped
		}
		if e.Err != nil {
			return nil, e.Err
		}
		return w.getServices()
	}
}

func (w *watcher) Stop() error {
	w.cancel()
	return nil
}

func (w *watcher) getServices() ([]*registry.ServiceInstance, error) {
	servicesID, _, err := w.conn.Children(w.prefix)
	if err != nil {
		return nil, err
	}
	items := make([]*registry.ServiceInstance, 0, len(servicesID))
	for _, id := range servicesID {
		servicePath := path.Join(w.prefix, id)
		b, _, err := w.conn.Get(servicePath)
		if err != nil {
			return nil, err
		}
		item, err := unmarshal(b)
		if err != nil {
			return nil, err
		}

		// 若实例服务名与 watcher 不一致则跳过
		if item.Name != w.serviceName {
			continue
		}

		items = append(items, item)
	}
	return items, nil
}
