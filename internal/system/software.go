package system

import (
	"os/exec"
	"strings"
)

type Software struct {
	Name        string `json:"name"`
	Version     string `json:"version,omitempty"`
	Status      string `json:"status"`
	Path        string `json:"path,omitempty"`
	Description string `json:"description,omitempty"`
}

// GetSoftwareList 返回系统中安装的软件列表
func GetSoftwareList() ([]Software, error) {
	// 这里是一个示例实现，你可以根据需要扩展
	// 例如：检查常见的包管理器、系统服务等

	// 获取系统服务列表
	cmd := exec.Command("systemctl", "list-units", "--type=service", "--no-pager")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	services := []Software{}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, ".service") {
			fields := strings.Fields(line)
			if len(fields) >= 4 {
				name := strings.TrimSuffix(fields[0], ".service")
				status := "stopped"
				if strings.Contains(line, "running") {
					status = "running"
				}

				services = append(services, Software{
					Name:   name,
					Status: status,
				})
			}
		}
	}

	return services, nil
}
