package managers

import "fmt"

// ProxyManager 代理软件管理相关功能
type ProxyManager struct {
	*SoftManager
}

// GetProxySoftwares 获取所有的代理软件
func (p *ProxyManager) GetProxySoftwares() []string {
	proxySoftwares := make([]string, 0)
	for name, software := range p.Softwares {
		if software.GetInfo().IsProxySoftware {
			proxySoftwares = append(proxySoftwares, name)
		}
	}
	return proxySoftwares
}

// IsProxyOn 判断是否开启了代理
func (p *ProxyManager) IsProxyOn() bool {
	return p.checkEnvProxy() || p.checkSystemProxy() || p.checkHTTPClientProxy()
}

// OpenProxy 打开代理
func (p *ProxyManager) OpenProxy() (string, error) {
	proxySoftwares := p.GetProxySoftwares()
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

// CloseProxy 关闭代理
func (p *ProxyManager) CloseProxy(softwareName string) error {
	soft, err := p.GetSoftware(softwareName)
	if err != nil {
		return err
	}
	return soft.Stop()
}
