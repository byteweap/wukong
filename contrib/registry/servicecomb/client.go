package servicecomb

import (
	"github.com/go-chassis/cari/discovery"
	"github.com/go-chassis/sc-client"
)

type RegistryClient interface {
	GetMicroServiceID(appID, microServiceName, version, env string, opts ...sc.CallOption) (string, error)
	FindMicroServiceInstances(consumerID, appID, microServiceName, versionRule string, opts ...sc.CallOption) ([]*discovery.MicroServiceInstance, error)
	RegisterService(microService *discovery.MicroService) (string, error)
	RegisterMicroServiceInstance(microServiceInstance *discovery.MicroServiceInstance) (string, error)
	Heartbeat(microServiceID, microServiceInstanceID string) (bool, error)
	UnregisterMicroServiceInstance(microServiceID, microServiceInstanceID string) (bool, error)
	WatchMicroService(microServiceID string, callback func(*sc.MicroServiceInstanceChangedEvent)) error
}
