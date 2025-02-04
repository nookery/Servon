package system

import (
	"fmt"
	"os/exec"
	"strings"

	"servon/internal/utils"
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
	utils.Info("checking status for service: %s", service)
	cmd := exec.Command("service", service, "status")
	output, err := cmd.CombinedOutput()
	status := strings.TrimSpace(string(output))

	if err != nil {
		utils.Error("status check for %s failed: %v (output: %s)", service, err, status)
		return false
	}

	isRunning := strings.Contains(status, "is running")
	utils.Info("[service] service %s status: %s", service, status)
	return isRunning
}

func (s ServiceCmd) Stop(service string) error {
	cmd := exec.Command("sudo", "service", service, "stop")
	if output, err := cmd.CombinedOutput(); err != nil {
		utils.Error("[service] stop service %s failed: %v\n%s", service, err, string(output))
		return fmt.Errorf("failed to stop service: %v\n%s", err, string(output))
	}
	return nil
}

func (s ServiceCmd) Reload(service string) error {
	cmd := exec.Command("sudo", "service", service, "reload")
	if output, err := cmd.CombinedOutput(); err != nil {
		utils.Error("[service] reload service %s failed: %v\n%s", service, err, string(output))
		return fmt.Errorf("failed to reload service: %v\n%s", err, string(output))
	}
	return nil
}

func (s ServiceCmd) Start(service string) error {
	cmd := exec.Command("sudo", "service", service, "start")
	if output, err := cmd.CombinedOutput(); err != nil {
		utils.Error("[service] start service %s failed: %v\n%s", service, err, string(output))
		return fmt.Errorf("failed to start service: %v\n%s", err, string(output))
	}
	return nil
}
