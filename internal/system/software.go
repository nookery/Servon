package system

import (
	"fmt"
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

// 支持安装的软件列表
var supportedSoftware = []Software{
	{
		Name:        "nginx",
		Description: "高性能的 HTTP 和反向代理服务器",
	},
	{
		Name:        "mysql",
		Description: "流行的关系型数据库",
	},
	{
		Name:        "redis",
		Description: "内存数据结构存储系统",
	},
	{
		Name:        "docker",
		Description: "应用容器引擎",
	},
	{
		Name:        "postgresql",
		Description: "开源对象关系数据库系统",
	},
}

// GetSoftwareList 返回系统中已安装和可安装的软件列表
func GetSoftwareList() ([]Software, error) {
	// 获取系统服务列表
	cmd := exec.Command("systemctl", "list-units", "--type=service", "--no-pager")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// 解析已安装的服务
	installedServices := make(map[string]string) // name -> status
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
				installedServices[name] = status
			}
		}
	}

	// 合并已安装和支持安装的软件列表
	result := []Software{}

	// 添加已安装的软件
	for _, sw := range supportedSoftware {
		if status, exists := installedServices[sw.Name]; exists {
			// 软件已安装
			sw.Status = status
			// 尝试获取版本信息
			sw.Version = getVersion(sw.Name)
			result = append(result, sw)
		} else {
			// 软件未安装
			sw.Status = "not_installed"
			result = append(result, sw)
		}
	}

	return result, nil
}

// getVersion 尝试获取软件版本
func getVersion(name string) string {
	var cmd *exec.Cmd
	switch name {
	case "nginx":
		cmd = exec.Command("nginx", "-v")
	case "mysql":
		cmd = exec.Command("mysql", "--version")
	case "redis":
		cmd = exec.Command("redis-server", "--version")
	case "docker":
		cmd = exec.Command("docker", "--version")
	case "postgresql":
		cmd = exec.Command("psql", "--version")
	default:
		return ""
	}

	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	// 简单返回第一行作为版本信息
	version := strings.Split(string(output), "\n")[0]
	return strings.TrimSpace(version)
}

// InstallSoftware 安装指定的软件
func InstallSoftware(name string) error {
	// 检查软件是否在支持列表中
	var targetSoftware *Software
	for _, sw := range supportedSoftware {
		if sw.Name == name {
			targetSoftware = &sw
			break
		}
	}
	if targetSoftware == nil {
		return fmt.Errorf("不支持安装该软件: %s", name)
	}

	// 根据不同的软件使用不同的安装命令
	var cmd *exec.Cmd
	switch name {
	case "nginx":
		cmd = exec.Command("apt-get", "install", "-y", "nginx")
	case "mysql":
		cmd = exec.Command("apt-get", "install", "-y", "mysql-server")
	case "redis":
		cmd = exec.Command("apt-get", "install", "-y", "redis-server")
	case "docker":
		cmd = exec.Command("apt-get", "install", "-y", "docker.io")
	case "postgresql":
		cmd = exec.Command("apt-get", "install", "-y", "postgresql")
	default:
		return fmt.Errorf("未知的软件: %s", name)
	}

	// 执行安装命令
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("安装失败: %v, 输出: %s", err, string(output))
	}

	// 启动服务
	startCmd := exec.Command("systemctl", "start", name)
	if err := startCmd.Run(); err != nil {
		return fmt.Errorf("启动服务失败: %v", err)
	}

	return nil
}

// UninstallSoftware 卸载指定的软件
func UninstallSoftware(name string) error {
	// 检查软件是否在支持列表中
	var targetSoftware *Software
	for _, sw := range supportedSoftware {
		if sw.Name == name {
			targetSoftware = &sw
			break
		}
	}
	if targetSoftware == nil {
		return fmt.Errorf("不支持卸载该软件: %s", name)
	}

	// 停止服务
	stopCmd := exec.Command("systemctl", "stop", name)
	if err := stopCmd.Run(); err != nil {
		return fmt.Errorf("停止服务失败: %v", err)
	}

	// 根据不同的软件使用不同的卸载命令
	var cmd *exec.Cmd
	switch name {
	case "nginx":
		cmd = exec.Command("apt-get", "remove", "-y", "nginx")
	case "mysql":
		cmd = exec.Command("apt-get", "remove", "-y", "mysql-server")
	case "redis":
		cmd = exec.Command("apt-get", "remove", "-y", "redis-server")
	case "docker":
		cmd = exec.Command("apt-get", "remove", "-y", "docker.io")
	case "postgresql":
		cmd = exec.Command("apt-get", "remove", "-y", "postgresql")
	default:
		return fmt.Errorf("未知的软件: %s", name)
	}

	// 执行卸载命令
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("卸载失败: %v, 输出: %s", err, string(output))
	}

	return nil
}
