package registry

import "context"

// Discovery is service discovery.
type Discovery interface {
	// ID return the implement id.
	ID() string
	// GetService return the service instances in memory according to the service name.
	GetService(ctx context.Context, serviceName string) ([]*ServiceInstance, error)
	// Watch creates a watcher according to the service name.
	Watch(ctx context.Context, serviceName string) (Watcher, error)
}
