package system

import (
	"fmt"
	"os/exec"
	"strings"
)

// ServiceCmd provides methods to interact with services using service command
type ServiceCmd struct{}

func NewServiceCmd() ServiceCmd {
	return ServiceCmd{}
}

// Type returns the type of service manager
func (s ServiceCmd) Type() string {
	return "service"
}

func (s ServiceCmd) IsActive(service string) bool {
	fmt.Printf("[service] checking status for service: %s\n", service)
	cmd := exec.Command("service", service, "status")
	output, err := cmd.CombinedOutput()
	status := strings.TrimSpace(string(output))

	if err != nil {
		fmt.Printf("[service] status check for %s failed: %v (output: %s)\n", service, err, status)
		return false
	}

	isRunning := strings.Contains(status, "is running")
	fmt.Printf("[service] service %s status: %s\n", service, status)
	return isRunning
}

func (s ServiceCmd) Stop(service string) error {
	cmd := exec.Command("sudo", "service", service, "stop")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to stop service: %v\n%s", err, string(output))
	}
	return nil
}

func (s ServiceCmd) Reload(service string) error {
	cmd := exec.Command("sudo", "service", service, "reload")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to reload service: %v\n%s", err, string(output))
	}
	return nil
}
