package system

import (
	"fmt"
	"os/exec"
	"strings"
)

// SystemCtl provides methods to interact with systemd services
type SystemCtl struct{}

// NewSystemCtl creates a new SystemCtl instance
func NewSystemCtl() SystemCtl {
	return SystemCtl{}
}

// Type returns the type of service manager
func (s SystemCtl) Type() string {
	return "systemctl"
}

// IsActive checks if a service is active
func (s SystemCtl) IsActive(service string) bool {
	fmt.Printf("[systemctl] checking status for service: %s\n", service)
	cmd := exec.Command("systemctl", "is-active", service)
	output, err := cmd.CombinedOutput()
	status := strings.TrimSpace(string(output))

	if err != nil {
		fmt.Printf("[systemctl] is-active %s failed: %v (status: %s)\n", service, err, status)
		return false
	}

	fmt.Printf("[systemctl] service %s status: %s\n", service, status)
	return status == "active"
}

// Stop stops a service
func (s SystemCtl) Stop(service string) error {
	cmd := exec.Command("sudo", "systemctl", "stop", service)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to stop service: %v\n%s", err, string(output))
	}
	return nil
}

// Reload reloads a service
func (s SystemCtl) Reload(service string) error {
	cmd := exec.Command("sudo", "systemctl", "reload", service)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to reload service: %v\n%s", err, string(output))
	}
	return nil
}

// Start starts a service
func (s SystemCtl) Start(service string) error {
	cmd := exec.Command("sudo", "systemctl", "start", service)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to start service: %v\n%s", err, string(output))
	}
	return nil
}
