package managers

import (
	"fmt"
	"net/http"
	"os"
	"servon/core/internal/contract"
	"servon/core/internal/utils"
	"strconv"
)

// SoftManager 基础软件管理功能
type SoftManager struct {
	Softwares map[string]contract.Software
	Gateways  map[string]contract.SuperGateway
	LogDir    string
	LogUtil   *utils.LogUtil
	*ProxySoftManager
	*GatewaySoftManager
}

func NewSoftManager(logDir string) *SoftManager {
	sm := &SoftManager{
		Softwares: make(map[string]contract.Software),
		Gateways:  make(map[string]contract.SuperGateway),
		LogDir:    logDir,
		LogUtil:   utils.NewTopicLogUtil(logDir, "soft"),
	}

	sm.ProxySoftManager = &ProxySoftManager{SoftManager: sm}
	sm.GatewaySoftManager = &GatewaySoftManager{SoftManager: sm}

	sm.LogUtil.Info("初始化软件管理器")
	return sm
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

		PrintListWithTitle("可用的软件", registeredSoftwares)
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

// RegisterSoftware 注册普通软件
func (s *SoftManager) RegisterSoftware(name string, software contract.Software) error {
	s.LogUtil.Info("注册软件: " + name)
	if _, exists := s.Softwares[name]; exists {
		return PrintAndReturnError(fmt.Sprintf("软件 %s 已注册", name))
	}
	s.Softwares[name] = software
	return nil
}

// RegisterGateway 注册网关软件
func (c *SoftManager) RegisterGateway(name string, gateway contract.SuperGateway) error {
	c.LogUtil.Info("注册网关软件: " + name)

	// 检查是否已注册为普通软件或网关
	if _, exists := c.Softwares[name]; exists {
		return PrintAndReturnError(fmt.Sprintf("软件 %s 已注册为普通软件", name))
	}
	if _, exists := c.Gateways[name]; exists {
		return PrintAndReturnError(fmt.Sprintf("软件 %s 已注册为网关软件", name))
	}

	// 注册为网关软件
	c.Gateways[name] = gateway
	// 同时注册为普通软件（因为网关也是软件）
	c.Softwares[name] = gateway
	return nil
}

// GetGateway 获取网关软件
func (c *SoftManager) GetGateway(name string) (contract.SuperGateway, error) {
	c.LogUtil.Info("获取网关软件: " + name)
	gateway, ok := c.Gateways[name]
	if !ok {
		return nil, PrintAndReturnError(fmt.Sprintf("网关软件 %s 未注册", name))
	}
	c.LogUtil.Success("获取网关软件成功: " + name)
	return gateway, nil
}

// GetAllGateways 获取所有网关软件
func (c *SoftManager) GetAllGateways() []string {
	c.LogUtil.Info("获取所有网关软件...")
	gatewayNames := make([]string, 0, len(c.Gateways))
	for name := range c.Gateways {
		gatewayNames = append(gatewayNames, name)
	}
	return gatewayNames
}

// IsGateway 判断软件是否为网关软件
func (c *SoftManager) IsGateway(name string) bool {
	_, ok := c.Gateways[name]
	return ok
}

// GetAllSoftware 获取所有软件
func (c *SoftManager) GetAllSoftware() []string {
	c.LogUtil.Info("获取所有软件...")
	softwareNames := make([]string, 0, len(c.Softwares))
	for name := range c.Softwares {
		softwareNames = append(softwareNames, name)
	}

	c.LogUtil.Success("获取所有软件成功，共 " + strconv.Itoa(len(softwareNames)) + " 个")
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
