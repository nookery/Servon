package soft

import (
	"fmt"
	"net/http"
	"os"
	"servon/core/internal/contract"
	"servon/core/internal/utils"
	"strconv"
)

// Manager 基础软件管理功能
type Manager struct {
	Softwares map[string]contract.Software
	Gateways  map[string]contract.SuperGateway
	LogDir    string
	LogUtil   *utils.LogUtil
	ShellUtil *utils.ShellUtil
	*ProxyManager
	*GatewayManager
	*AptManager
	*DpkgManager
}

// NewManager 创建新的软件管理器
func NewManager(logDir string) *Manager {
	sm := &Manager{
		Softwares: make(map[string]contract.Software),
		Gateways:  make(map[string]contract.SuperGateway),
		LogDir:    logDir,
		LogUtil:   utils.NewTopicLogUtil(logDir, "soft"),
		ShellUtil: utils.NewShellUtil(),
	}

	sm.ProxyManager = &ProxyManager{Manager: sm}
	sm.GatewayManager = &GatewayManager{Manager: sm}
	sm.AptManager = &AptManager{Manager: sm}
	sm.DpkgManager = &DpkgManager{Manager: sm}
	sm.LogUtil.Info("初始化软件管理器")
	return sm
}

// GetSoftware 获取软件
func (p *Manager) GetSoftware(name string) (contract.SuperSoft, error) {
	software, ok := p.Softwares[name]
	if !ok {
		return nil, fmt.Errorf("软件 %s 未注册", name)
	}
	return software, nil
}

// HasSoftware 判断软件是否存在
func (p *Manager) HasSoftware(name string) bool {
	_, ok := p.Softwares[name]
	return ok
}

// Install 安装软件
func (c *Manager) Install(name string) error {
	software, ok := c.Softwares[name]
	if !ok {
		registeredSoftwares := make([]string, 0, len(c.Softwares))
		for name := range c.Softwares {
			registeredSoftwares = append(registeredSoftwares, name)
		}

		return fmt.Errorf("软件 %s 未注册，可用的软件: %v", name, registeredSoftwares)
	}

	return software.Install()
}

// UninstallSoftware 卸载软件
func (c *Manager) UninstallSoftware(name string) error {
	software, ok := c.Softwares[name]
	if !ok {
		return fmt.Errorf("软件 %s 未注册", name)
	}
	return software.Uninstall()
}

// StartSoftware 启动软件
func (c *Manager) StartSoftware(name string) error {
	software, ok := c.Softwares[name]
	if !ok {
		return fmt.Errorf("软件 %s 未注册", name)
	}
	return software.Start()
}

// StopSoftware 停止软件
func (c *Manager) StopSoftware(name string) error {
	software, ok := c.Softwares[name]
	if !ok {
		return fmt.Errorf("软件 %s 未注册", name)
	}
	return software.Stop()
}

// GetSoftwareStatus 获取软件状态
func (c *Manager) GetSoftwareStatus(name string) (map[string]string, error) {
	software, ok := c.Softwares[name]
	if !ok {
		return nil, fmt.Errorf("软件 %s 未注册", name)
	}

	return software.GetStatus()
}

// RegisterSoftware 注册普通软件
func (s *Manager) RegisterSoftware(name string, software contract.Software) error {
	if _, exists := s.Softwares[name]; exists {
		return fmt.Errorf("软件 %s 已注册", name)
	}
	s.Softwares[name] = software
	return nil
}

// GetAllSoftware 获取所有软件
func (c *Manager) GetAllSoftware() []string {
	softwareNames := make([]string, 0, len(c.Softwares))
	for name := range c.Softwares {
		softwareNames = append(softwareNames, name)
	}

	c.LogUtil.Success("获取所有软件成功，共 " + strconv.Itoa(len(softwareNames)) + " 个")
	return softwareNames
}

// 检查环境变量中的代理设置
func (p *Manager) checkEnvProxy() bool {
	proxyEnvVars := []string{
		"http_proxy",
		"HTTP_PROXY",
		"https_proxy",
		"HTTPS_PROXY",
		"all_proxy",
		"ALL_PROXY",
	}

	for _, env := range proxyEnvVars {
		if proxy := os.Getenv(env); proxy != "" {
			return true
		}
	}
	return false
}

// 检查系统代理设置
func (p *Manager) checkSystemProxy() bool {
	// 创建一个完整的请求对象，包含必要的URL
	req, err := http.NewRequest("GET", "https://github.com", nil)
	if err != nil {
		return false
	}

	// 获取系统代理设置
	proxyURL, err := http.ProxyFromEnvironment(req)
	if err != nil {
		return false
	}
	return proxyURL != nil
}

// 检查默认 HTTP 客户端的代理设置
func (p *Manager) checkHTTPClientProxy() bool {
	transport, ok := http.DefaultClient.Transport.(*http.Transport)
	if !ok {
		return false
	}

	if transport.Proxy != nil {
		// 尝试获取一个示例请求的代理
		req, _ := http.NewRequest("GET", "https://github.com", nil)
		proxyURL, err := transport.Proxy(req)
		if err == nil && proxyURL != nil {
			return true
		}
	}

	return false
}
