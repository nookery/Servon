package system

// ServiceManager defines the interface for service management
type ServiceManager interface {
	Type() string
	IsActive(service string) bool
	Stop(service string) error
	Start(service string) error
	Reload(service string) error
}
