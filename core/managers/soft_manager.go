package managers

import (
	"fmt"
	"net/http"
	"os"
	"servon/components/shell_util"
	"servon/components/soft_util"

	"servon/core/contract"
)

// SoftManager 基础软件管理功能
type SoftManager struct {
	Softwares map[string]contract.Software
	Gateways  map[string]contract.SuperGateway
	Services  map[string]contract.SuperService
	ShellUtil *shell_util.ShellUtil
	*ProxyManager
	*GatewayManager
	*soft_util.AptManager
	*soft_util.DpkgManager
	*ServiceSoftManager
}

// NewManager 创建新的软件管理器
func NewSoftManager() *SoftManager {
	sm := &SoftManager{
		Softwares: make(map[string]contract.Software),
		Gateways:  make(map[string]contract.SuperGateway),
		ShellUtil: shell_util.NewShellUtil(),
	}

	sm.ProxyManager = &ProxyManager{SoftManager: sm}
	sm.GatewayManager = &GatewayManager{SoftManager: sm}
	sm.AptManager = &soft_util.AptManager{}
	sm.DpkgManager = &soft_util.DpkgManager{}
	sm.ServiceSoftManager = &ServiceSoftManager{SoftManager: sm}
	return sm
}

// GetSoftware 获取软件
func (p *SoftManager) GetSoftware(name string) (contract.SuperSoft, error) {
	software, ok := p.Softwares[name]
	if !ok {
		return nil, fmt.Errorf("软件 %s 未注册", name)
	}
	return software, nil
}

// HasSoftware 判断软件是否存在
func (p *SoftManager) HasSoftware(name string) bool {
	_, ok := p.Softwares[name]
	return ok
}

// Install 安装软件
func (c *SoftManager) Install(name string) error {
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
func (c *SoftManager) UninstallSoftware(name string) error {
	software, ok := c.Softwares[name]
	if !ok {
		return fmt.Errorf("软件 %s 未注册", name)
	}
	return software.Uninstall()
}

// StartSoftware 启动软件
func (c *SoftManager) StartSoftware(name string) error {
	software, ok := c.Softwares[name]
	if !ok {
		return fmt.Errorf("软件 %s 未注册", name)
	}
	return software.Start()
}

// StopSoftware 停止软件
func (c *SoftManager) StopSoftware(name string) error {
	software, ok := c.Softwares[name]
	if !ok {
		return fmt.Errorf("软件 %s 未注册", name)
	}
	return software.Stop()
}

// GetSoftwareStatus 获取软件状态
func (c *SoftManager) GetSoftwareStatus(name string) (map[string]string, error) {
	software, ok := c.Softwares[name]
	if !ok {
		return nil, fmt.Errorf("软件 %s 未注册", name)
	}

	return software.GetStatus()
}

// RegisterSoftware 注册普通软件
func (s *SoftManager) RegisterSoftware(name string, software contract.Software) error {
	if _, exists := s.Softwares[name]; exists {
		return fmt.Errorf("软件 %s 已注册", name)
	}
	s.Softwares[name] = software
	return nil
}

// GetAllSoftware 获取所有软件
func (c *SoftManager) GetAllSoftware() []string {
	softwareNames := make([]string, 0, len(c.Softwares))
	for name := range c.Softwares {
		softwareNames = append(softwareNames, name)
	}

	return softwareNames
}

// 检查环境变量中的代理设置
func (p *SoftManager) checkEnvProxy() bool {
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
func (p *SoftManager) checkSystemProxy() bool {
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
func (p *SoftManager) checkHTTPClientProxy() bool {
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
