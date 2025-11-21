package registry

import "context"

// Registrar is service registrar.
type Registrar interface {
	// ID return the implement id.
	ID() string
	// Register the registration.
	Register(ctx context.Context, service *ServiceInstance) error
	// Deregister the registration.
	Deregister(ctx context.Context, service *ServiceInstance) error
}
