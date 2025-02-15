package managers

import (
	"fmt"
	"net/http"
	"os"
	"servon/core/internal/contract"
)

type SoftManager struct {
	Softwares map[string]contract.SuperSoft
}

func NewSoftManager() *SoftManager {
	return &SoftManager{
		Softwares: make(map[string]contract.SuperSoft),
	}
}

// GetProxySoftwares 获取所有的代理软件
func (p *SoftManager) GetProxySoftwares() []string {
	proxySoftwares := make([]string, 0)
	for name, software := range p.Softwares {
		if software.GetInfo().IsProxySoftware {
			proxySoftwares = append(proxySoftwares, name)
		}
	}
	return proxySoftwares
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

// Install 安装软件, 如果提供了日志通道则输出日志
func (c *SoftManager) Install(name string) error {
	software, ok := c.Softwares[name]
	if !ok {
		registeredSoftwares := make([]string, 0, len(c.Softwares))
		for name := range c.Softwares {
			registeredSoftwares = append(registeredSoftwares, name)
		}

		PrintList(registeredSoftwares, "可用的软件")
		return PrintAndReturnError(fmt.Sprintf("软件 %s 未注册", name))
	}

	return software.Install()
}

// UninstallSoftware 卸载软件
func (c *SoftManager) UninstallSoftware(name string) error {
	software, ok := c.Softwares[name]
	if !ok {
		return PrintAndReturnError(fmt.Sprintf("软件 %s 未注册", name))
	}
	return software.Uninstall()
}

// StartSoftware 启动软件
func (c *SoftManager) StartSoftware(name string) error {
	software, ok := c.Softwares[name]
	if !ok {
		return PrintAndReturnError(fmt.Sprintf("软件 %s 未注册", name))
	}
	return software.Start()
}

// StopSoftware 停止软件
func (c *SoftManager) StopSoftware(name string) error {
	software, ok := c.Softwares[name]
	if !ok {
		return PrintAndReturnError(fmt.Sprintf("软件 %s 未注册", name))
	}
	return software.Stop()
}

// GetSoftwareStatus 获取软件状态
func (c *SoftManager) GetSoftwareStatus(name string) (map[string]string, error) {
	software, ok := c.Softwares[name]
	if !ok {
		return nil, PrintAndReturnError(fmt.Sprintf("软件 %s 未注册", name))
	}

	return software.GetStatus()
}

// RegisterSoftware 注册软件
func (c *SoftManager) RegisterSoftware(name string, software contract.SuperSoft) error {
	if _, exists := c.Softwares[name]; exists {
		return PrintAndReturnError(fmt.Sprintf("软件 %s 已注册", name))
	}
	c.Softwares[name] = software
	return nil
}

// GetAllSoftware 获取所有软件
func (c *SoftManager) GetAllSoftware() []string {
	PrintInfo("获取所有软件...")
	softwareNames := make([]string, 0, len(c.Softwares))
	for name := range c.Softwares {
		softwareNames = append(softwareNames, name)
	}
	return softwareNames
}

// IsProxyOn 判断是否开启了代理
func (p *SoftManager) IsProxyOn() bool {
	// 方法1: 检查环境变量
	if p.checkEnvProxy() {
		return true
	}

	// 方法2: 检查系统代理设置
	if p.checkSystemProxy() {
		return true
	}

	// 方法3: 检查 HTTP_PROXY 客户端设置
	if p.checkHTTPClientProxy() {
		return true
	}

	return false
}

// OpenProxy 打开代理，打开成功就返回代理软件的名字，否则返回错误
func (p *SoftManager) OpenProxy() (string, error) {
	// 获取所有的代理软件
	proxySoftwares := p.GetProxySoftwares()

	// 尝试逐个启动，如果有启动成功，就返回
	for _, software := range proxySoftwares {
		soft, err := p.GetSoftware(software)
		if err != nil {
			return "", err
		}

		err = soft.Start()
		if err != nil {
			return "", err
		}

		return software, nil
	}

	return "", fmt.Errorf("没有找到可用的代理软件")
}

// CloseProxy 关闭代理，配合OpenProxy使用
func (p *SoftManager) CloseProxy(softwareName string) error {
	soft, err := p.GetSoftware(softwareName)
	if err != nil {
		return err
	}

	err = soft.Stop()
	if err != nil {
		return err
	}

	return nil
}

// checkEnvProxy 检查环境变量中的代理设置
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

// checkSystemProxy 检查系统代理设置
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

// checkHTTPClientProxy 检查默认 HTTP 客户端的代理设置
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
