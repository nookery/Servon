package system

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"sync"
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
	// 获取已安装的软件包
	cmd := exec.Command("dpkg", "-l")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// 解析已安装的软件包
	installedPackages := make(map[string]bool)
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ii") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				installedPackages[fields[1]] = true
			}
		}
	}

	// 获取系统服务状态
	serviceCmd := exec.Command("systemctl", "list-units", "--type=service", "--all", "--no-pager")
	serviceOutput, err := serviceCmd.Output()
	if err != nil {
		return nil, err
	}

	// 解析服务状态
	serviceStatus := make(map[string]string)
	serviceLines := strings.Split(string(serviceOutput), "\n")
	for _, line := range serviceLines {
		if strings.Contains(line, ".service") {
			fields := strings.Fields(line)
			if len(fields) >= 4 {
				name := strings.TrimSuffix(fields[0], ".service")
				status := "stopped"
				if strings.Contains(line, "running") {
					status = "running"
				}
				serviceStatus[name] = status
			}
		}
	}

	// 合并已安装和支持安装的软件列表
	result := []Software{}

	// 添加软件
	for _, sw := range supportedSoftware {
		software := sw

		// 检查软件包是否已安装
		isInstalled := false
		switch sw.Name {
		case "nginx":
			isInstalled = installedPackages["nginx"]
		case "mysql":
			isInstalled = installedPackages["mysql-server"]
		case "redis":
			isInstalled = installedPackages["redis-server"]
		case "docker":
			isInstalled = installedPackages["docker.io"]
		case "postgresql":
			isInstalled = installedPackages["postgresql"]
		}

		if isInstalled {
			// 软件已安装，检查服务状态
			if status, exists := serviceStatus[sw.Name]; exists {
				software.Status = status
			} else {
				software.Status = "stopped"
			}
			// 获取版本信息
			software.Version = getVersion(sw.Name)
		} else {
			software.Status = "not_installed"
		}

		result = append(result, software)
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

// InstallSoftware 安装指定的软件，并返回命令输出的通道
func InstallSoftware(name string) (chan string, error) {
	// 创建输出通道
	outputChan := make(chan string, 100)

	// 检查软件是否在支持列表中
	var targetSoftware *Software
	for _, sw := range supportedSoftware {
		if sw.Name == name {
			targetSoftware = &sw
			break
		}
	}
	if targetSoftware == nil {
		close(outputChan)
		return outputChan, fmt.Errorf("不支持安装该软件: %s", name)
	}

	// 发送开始安装的消息
	outputChan <- fmt.Sprintf("开始安装 %s...", name)

	// 获取当前用户信息
	whoamiCmd := exec.Command("whoami")
	whoamiOutput, err := whoamiCmd.Output()
	if err != nil {
		outputChan <- fmt.Sprintf("获取用户信息失败: %v", err)
	} else {
		username := strings.TrimSpace(string(whoamiOutput))
		outputChan <- fmt.Sprintf("当前用户: %s", username)
	}

	// 获取用户组信息
	groupsCmd := exec.Command("groups")
	groupsOutput, err := groupsCmd.Output()
	if err != nil {
		outputChan <- fmt.Sprintf("获取用户组信息失败: %v", err)
	} else {
		groups := strings.TrimSpace(string(groupsOutput))
		outputChan <- fmt.Sprintf("用户组: %s", groups)
	}

	// 先尝试更新包索引
	outputChan <- "更新软件包索引..."
	updateCmd := exec.Command("sudo", "apt-get", "update")
	updateOutput, err := updateCmd.CombinedOutput()
	if err != nil {
		outputChan <- fmt.Sprintf("更新软件包索引失败: %v\n%s", err, string(updateOutput))
	} else {
		outputChan <- "软件包索引更新完成"
	}

	// 根据不同的软件使用不同的安装命令
	var cmd *exec.Cmd
	switch name {
	case "nginx":
		cmd = exec.Command("sudo", "apt-get", "install", "-y", "nginx")
	case "mysql":
		cmd = exec.Command("sudo", "apt-get", "install", "-y", "mysql-server")
	case "redis":
		cmd = exec.Command("sudo", "apt-get", "install", "-y", "redis-server")
	case "docker":
		cmd = exec.Command("sudo", "apt-get", "install", "-y", "docker.io")
	case "postgresql":
		cmd = exec.Command("sudo", "apt-get", "install", "-y", "postgresql")
	default:
		close(outputChan)
		return outputChan, fmt.Errorf("未知的软件: %s", name)
	}

	// 获取命令的输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		close(outputChan)
		return outputChan, fmt.Errorf("创建输出管道失败: %v", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		close(outputChan)
		return outputChan, fmt.Errorf("创建错误输出管道失败: %v", err)
	}

	// 启动命令
	outputChan <- fmt.Sprintf("执行安装命令: %s", cmd.String())
	if err := cmd.Start(); err != nil {
		close(outputChan)
		return outputChan, fmt.Errorf("启动安装命令失败: %v", err)
	}

	// 在后台处理输出
	go func() {
		defer close(outputChan)

		// 使用 WaitGroup 等待所有输出处理完成
		var wg sync.WaitGroup
		wg.Add(2)

		// 读取标准输出
		go func() {
			defer wg.Done()
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				outputChan <- scanner.Text()
			}
			if err := scanner.Err(); err != nil {
				outputChan <- fmt.Sprintf("读取标准输出错误: %v", err)
			}
		}()

		// 读取错误输出
		go func() {
			defer wg.Done()
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				outputChan <- fmt.Sprintf("Error: %s", scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				outputChan <- fmt.Sprintf("读取错误输出错误: %v", err)
			}
		}()

		// 等待所有输出读取完成
		wg.Wait()

		// 等待命令完成
		if err := cmd.Wait(); err != nil {
			outputChan <- fmt.Sprintf("安装过程出错: %v", err)
			return
		}

		outputChan <- "软件安装完成，正在启动服务..."

		// 启动服务
		startCmd := exec.Command("sudo", "systemctl", "start", name)
		if output, err := startCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("启动服务失败: %v\n%s", err, string(output))
			return
		}

		outputChan <- "安装完成"
	}()

	return outputChan, nil
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
