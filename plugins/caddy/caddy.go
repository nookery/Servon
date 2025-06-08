package caddy

import (
	"fmt"
	"os"
	"os/exec"
	"servon/core"
	"strings"
)

type Caddy struct {
	BaseDir string
	CaddyTemplate
	*core.App
	info core.SoftwareInfo
}

func NewCaddy(app *core.App) *Caddy {
	return &Caddy{
		App: app,
		info: core.SoftwareInfo{
			Name:            "caddy",
			Description:     "Modern web server with automatic HTTPS",
			IsProxySoftware: true,
			IsGateway:       true,
		},
	}
}

// GetInfo 获取软件信息
func (c *Caddy) GetInfo() core.SoftwareInfo {
	return c.info
}

// Install 安装 Caddy
func (c *Caddy) Install() error {
	osType := c.GetOSType()

	switch osType {
	case core.Ubuntu, core.Debian:
		// 注意：apt-key 已弃用，但在某些环境中仍然是最可靠的方法
		aptKeyCmd := `
		# 直接使用apt-key添加
		curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | apt-key add -
		
		# 添加仓库源
		echo "deb https://dl.cloudsmith.io/public/caddy/stable/deb/debian any-version main" > /etc/apt/sources.list.d/caddy-stable.list
		`
		err, _ := c.RunShell("sh", "-c", aptKeyCmd)
		if err != nil {
			unsafeCmd := `
			# 添加允许非安全仓库的配置
			mkdir -p /etc/apt/apt.conf.d
			echo 'APT::Get::AllowUnauthenticated "true";' > /etc/apt/apt.conf.d/99allow-unauth
			echo 'Acquire::AllowInsecureRepositories "true";' >> /etc/apt/apt.conf.d/99allow-unauth
			echo 'Acquire::AllowDowngradeToInsecureRepositories "true";' >> /etc/apt/apt.conf.d/99allow-unauth
			
			# 添加仓库源（不带签名验证）
			echo "deb [trusted=yes] https://dl.cloudsmith.io/public/caddy/stable/deb/debian any-version main" > /etc/apt/sources.list.d/caddy-stable.list
			`
			err, unsafeOutput := c.RunShell("sh", "-c", unsafeCmd)
			fmt.Printf("非安全方式输出: \n%s\n", unsafeOutput)
			if err != nil {
				fmt.Printf("非安全方式失败: %v\n", err)
				return err
			}
		}

		// 更新软件包索引
		fmt.Println("更新软件包索引...")
		output, err := c.AptUpdate()
		fmt.Printf("更新软件包索引输出: \n%s\n", output)
		if err != nil {
			fmt.Printf("更新软件包索引失败: %v\n", err)
			return err
		}

		// 安装 Caddy
		fmt.Println("安装 Caddy...")
		if err := c.AptInstall("caddy"); err != nil {
			fmt.Printf("Caddy 安装失败: %v\n", err)
			return err
		}

		fmt.Println("Caddy 安装完成")
	case core.CentOS, core.RedHat:
		errMsg := "暂不支持在 RHEL 系统上安装 Caddy"
		fmt.Printf("%s\n", errMsg)
		return fmt.Errorf("%s", errMsg)

	default:
		errMsg := fmt.Sprintf("不支持的操作系统: %s", osType)
		fmt.Printf("%s\n", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	// 验证安装结果
	if !c.IsInstalled("caddy") {
		errMsg := "Caddy 安装验证失败，未检测到已安装的包"
		fmt.Printf("%s\n", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	fmt.Println("Caddy 安装完成")

	return nil
}

// Uninstall 卸载 Caddy
func (c *Caddy) Uninstall() error {
	// 卸载软件包及其依赖
	fmt.Println("卸载软件包及其依赖...")
	if err := c.AptRemove("caddy"); err != nil {
		fmt.Printf("卸载软件包及其依赖失败:\n%s\n", err)
		return fmt.Errorf("卸载软件包及其依赖失败:\n%s", err)
	}

	// 清理配置文件
	if err := c.AptPurge("caddy"); err != nil {
		fmt.Printf("清理配置文件失败:\n%s\n", err)
		return fmt.Errorf("清理配置文件失败:\n%s", err)
	}

	// 删除源文件
	err, rmOutput := c.RunShell("rm", "/etc/apt/sources.list.d/caddy-stable.list")
	if err != nil {
		fmt.Printf("删除源文件失败:\n%s\n", err)
		return fmt.Errorf("删除源文件失败:\n%s", err)
	}
	fmt.Printf("删除源文件输出: %s\n", rmOutput)

	// 清理自动安装的依赖
	err, cleanOutput := c.RunShell("apt-get", "autoremove", "-y")
	if err != nil {
		fmt.Printf("清理依赖失败:\n%s\n", err)
		return fmt.Errorf("清理依赖失败:\n%s", err)
	}
	fmt.Printf("清理依赖输出: %s\n", cleanOutput)

	fmt.Println("Caddy 卸载完成")

	return nil
}

func (c *Caddy) GetStatus() (map[string]string, error) {
	if !c.IsInstalled("caddy") {
		return map[string]string{
			"status":  "not_installed",
			"version": "",
		}, nil
	}

	status := "stopped"
	if c.IsActive("caddy") {
		status = "running"
	}

	// 获取版本
	version := ""
	err, verOutput := c.RunShell("caddy", "version")
	if err == nil {
		version = strings.TrimSpace(verOutput)
	}

	return map[string]string{
		"status":  status,
		"version": version,
	}, nil
}

func (c *Caddy) Stop() error {
	err, output := c.RunShell("caddy", "stop")
	fmt.Printf("停止Caddy输出: %s\n", output)
	return err
}

// Reload reloads the Caddy configuration
func (c *Caddy) Reload() error {
	// 输出配置文件的存储目录
	fmt.Printf("配置文件的存储目录: %s\n", c.GetConfigDir())

	// 检查 caddy 是否在运行
	running, err := c.isRunning()
	if err != nil {
		fmt.Println("检查 caddy 运行状态失败")
		return fmt.Errorf("检查 caddy 运行状态失败")
	}

	if !running {
		fmt.Println("Caddy 服务未运行，请先启动 Caddy")
		return fmt.Errorf("Caddy 服务未运行，请先启动 Caddy")
	}

	err, output := c.RunShell("caddy", "reload", "--config", c.GetCaddyfilePath())
	fmt.Printf("重载Caddy输出: %s\n", output)
	return err
}

// isRunning 检查 caddy 是否在运行
func (c *Caddy) isRunning() (bool, error) {
	err, _ := c.RunShell("pgrep", "caddy")
	if err != nil {
		// exit status 1 表示没有找到进程
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Start starts the Caddy service
func (c *Caddy) Start() error {
	// 检查是否已安装
	if !c.IsInstalled("caddy") {
		errMsg := "Caddy 未安装，请先安装"
		fmt.Printf("%s\n", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	// 获取当前状态
	status, err := c.GetStatus()
	if err != nil {
		errMsg := fmt.Sprintf("获取状态失败: %v", err)
		fmt.Printf("%s\n", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	// 如果已经在运行，则不需要启动
	if status["status"] == "running" {
		fmt.Println("Caddy 服务已在运行中")
		return nil
	}

	fmt.Println("正在启动 Caddy 服务...")

	// 确保配置目录和 Caddyfile 存在
	if err := c.EnsureConfigDir(); err != nil {
		errMsg := fmt.Sprintf("创建 Caddy 配置目录失败: %v", err)
		fmt.Printf("%s\n", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	if err := c.EnsureCaddyfile(); err != nil {
		errMsg := fmt.Sprintf("确保 Caddyfile 存在失败: %v", err)
		fmt.Printf("%s\n", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	// 使用 RunShell 来启动 Caddy
	err, output := c.RunShell("caddy", "start", "--config", c.GetCaddyfilePath())
	fmt.Printf("启动Caddy输出: %s\n", output)
	if err != nil {
		errMsg := fmt.Sprintf("启动 Caddy 失败: %v", err)
		fmt.Printf("%s\n", errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	fmt.Println("Caddy 服务已成功启动")
	return nil
}

// Gateway 接口实现
func (c *Caddy) GetConfig() (map[string]interface{}, error) {
	// 确保配置目录和文件存在
	if err := c.EnsureConfigDir(); err != nil {
		return nil, fmt.Errorf("确保配置目录存在失败: %v", err)
	}

	if err := c.EnsureCaddyfile(); err != nil {
		return nil, fmt.Errorf("确保配置文件存在失败: %v", err)
	}

	content, err := os.ReadFile(c.GetCaddyfilePath())
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"config": string(content),
	}, nil
}

func (c *Caddy) SetConfig(config map[string]interface{}) error {
	if configStr, ok := config["config"].(string); ok {
		err := os.WriteFile(c.GetCaddyfilePath(), []byte(configStr), 0644)
		if err != nil {
			return err
		}
		return c.Reload()
	}
	return fmt.Errorf("invalid config format")
}

func (c *Caddy) GetProjects() ([]core.Project, error) {
	// TODO: 解析 Caddyfile 获取项目列表
	return []core.Project{}, nil
}

func (c *Caddy) AddProject(project core.Project) error {
	// TODO: 将项目配置添加到 Caddyfile
	return c.Reload()
}

func (c *Caddy) RemoveProject(projectName string) error {
	// TODO: 从 Caddyfile 中移除项目配置
	return c.Reload()
}

func (c *Caddy) ReloadConfig() error {
	return c.Reload()
}

// AddProxy 添加反向代理配置
// domain: 要代理的域名
// target: 目标地址（例如：127.0.0.1:8888）
func (c *Caddy) AddProxy(domain string, target string) error {
	// 确保参数合理
	if domain == "" || target == "" {
		return fmt.Errorf("域名和目标地址不能为空")
	}

	// 确保目标地址格式正确
	if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		fmt.Println("目标地址格式不正确，必须以 http:// 或 https:// 开头")
		return fmt.Errorf("目标地址格式不正确，必须以 http:// 或 https:// 开头")
	}

	// 确保配置目录和 Caddyfile 存在
	if err := c.EnsureConfigDir(); err != nil {
		return fmt.Errorf("创建配置目录失败: %v", err)
	}

	proxyConfig, err := c.RenderProxyConfig(domain, target)
	if err != nil {
		return fmt.Errorf("渲染代理配置失败: %v", err)
	}

	// 使用域名作为配置文件名获取项目配置路径
	configPath := c.GetProjectConfigPath(domain)

	// 写入代理配置
	if err := c.WriteConfig(configPath, proxyConfig); err != nil {
		return fmt.Errorf("写入代理配置失败: %v", err)
	}

	fmt.Println("添加反向代理配置成功")
	fmt.Printf("代理配置文件: %s\n", configPath)

	// 重新加载 Caddy 使配置生效
	if err := c.Reload(); err != nil {
		fmt.Println("重新加载 Caddy 失败")
		return fmt.Errorf("重新加载 Caddy 失败")
	}

	return nil
}

// RemoveProxy 移除指定域名的代理配置
func (c *Caddy) RemoveProxy(domain string) error {
	configPath := c.GetProjectConfigPath(domain)

	// 删除代理配置文件
	if err := c.RemoveConfig(configPath); err != nil {
		return fmt.Errorf("删除代理配置失败: %v", err)
	}

	// 重新加载 Caddy 使配置生效
	if err := c.Reload(); err != nil {
		return fmt.Errorf("移除代理后重新加载 Caddy 失败: %v", err)
	}

	return nil
}
