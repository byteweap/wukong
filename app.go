package wukong

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"

	"github.com/byteweap/wukong/component/log"
	"github.com/byteweap/wukong/component/registry"
	"golang.org/x/sync/errgroup"
)

type AppInfo interface {
	ID() string
	Name() string
	Version() string
	Metadata() map[string]string
	Endpoint() []string
}

type App struct {
	opts     *options
	ctx      context.Context
	cancel   context.CancelFunc
	mu       sync.Mutex
	instance *registry.ServiceInstance
}

func New(opts ...Option) *App {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	if o.logger != nil {
		log.SetLogger(o.logger)
	}
	ctx, cancel := context.WithCancel(o.ctx)
	return &App{
		opts:   o,
		ctx:    ctx,
		cancel: cancel,
	}
}

// ID returns app instance id.
func (a *App) ID() string { return a.opts.id }

// Name returns service name.
func (a *App) Name() string { return a.opts.name }

// Version returns app version.
func (a *App) Version() string { return a.opts.version }

// Metadata returns service metadata.
func (a *App) Metadata() map[string]string { return a.opts.metadata }

// Endpoint returns endpoints.
func (a *App) Endpoint() []string {
	if a.instance != nil {
		return a.instance.Endpoints
	}
	return nil
}

// Run executes all OnStart hooks registered with the application's Lifecycle.
func (a *App) Run() error {

	if err := a.buildInstance(); err != nil {
		return err
	}

	sctx := NewContext(a.ctx, a)
	eg, ctx := errgroup.WithContext(sctx)
	wg := sync.WaitGroup{}

	for _, fn := range a.opts.beforeStart {
		if err := fn(sctx); err != nil {
			return err
		}
	}
	octx := NewContext(a.opts.ctx, a)
	for _, srv := range a.opts.servers {
		server := srv
		eg.Go(func() error {
			<-ctx.Done() // wait for stop signal
			stopCtx := context.WithoutCancel(octx)
			if a.opts.stopTimeout > 0 {
				var cancel context.CancelFunc
				stopCtx, cancel = context.WithTimeout(stopCtx, a.opts.stopTimeout)
				defer cancel()
			}
			return server.Stop(stopCtx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done() // here is to ensure server start has begun running before register, so defer is not needed
			return server.Start(octx)
		})
	}
	wg.Wait()

	if err := a.register(ctx); err != nil {
		return err
	}

	for _, fn := range a.opts.afterStart {
		if err := fn(sctx); err != nil {
			return err
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.sigs...)
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return nil
		case <-c:
			return a.Stop()
		}
	})
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	var err error
	for _, fn := range a.opts.afterStop {
		err = fn(sctx)
	}
	return err
}

// Stop gracefully stops the application.
func (a *App) Stop() (err error) {
	sctx := NewContext(a.ctx, a)
	for _, fn := range a.opts.beforeStop {
		err = fn(sctx)
	}

	if err = a.deregister(NewContext(a.ctx, a)); err != nil {
		return err
	}

	if a.cancel != nil {
		a.cancel()
	}
	return err
}

// buildInstance builds a service instance.
func (a *App) buildInstance() error {

	endpoints := make([]string, 0, len(a.opts.endpoints))
	for _, e := range a.opts.endpoints {
		endpoints = append(endpoints, e.String())
	}
	if len(endpoints) == 0 {
		for _, srv := range a.opts.servers {
			e, err := srv.Endpoint()
			if err != nil {
				return err
			}
			endpoints = append(endpoints, e.String())
		}
	}
	instance := &registry.ServiceInstance{
		ID:        a.opts.id,
		Name:      a.opts.name,
		Version:   a.opts.version,
		Metadata:  a.opts.metadata,
		Endpoints: endpoints,
	}
	a.mu.Lock()
	a.instance = instance
	a.mu.Unlock()

	return nil
}

// register registers the service with the registry.
func (a *App) register(ctx context.Context) error {

	a.mu.Lock()
	instance := a.instance
	a.mu.Unlock()

	if a.opts.registry != nil {
		regCtx, regCancel := context.WithTimeout(ctx, a.opts.registryTimeout)
		defer regCancel()
		if err := a.opts.registry.Register(regCtx, instance); err != nil {
			return err
		}
	}
	return nil
}

// deregister de-registers the service with the registry.
func (a *App) deregister(ctx context.Context) error {

	a.mu.Lock()
	instance := a.instance
	a.mu.Unlock()

	if a.opts.registry != nil && instance != nil {
		regCtx, regCancel := context.WithTimeout(ctx, a.opts.registryTimeout)
		defer regCancel()
		if err := a.opts.registry.Deregister(regCtx, instance); err != nil {
			return err
		}
	}
	return nil
}

type appKey struct{}

// NewContext returns a new Context that carries value.
func NewContext(ctx context.Context, s AppInfo) context.Context {
	return context.WithValue(ctx, appKey{}, s)
}

// FromContext returns the Transport value stored in ctx, if any.
func FromContext(ctx context.Context) (s AppInfo, ok bool) {
	s, ok = ctx.Value(appKey{}).(AppInfo)
	return
}
