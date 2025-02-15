package libs

import (
	"fmt"
	"net/http"
	"os"
)

type ProxyManager struct{}

func NewProxyManager() *ProxyManager {
	return &ProxyManager{}
}

// IsProxyOn 判断是否开启了代理
func (p *ProxyManager) IsProxyOn() bool {
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
func (p *ProxyManager) OpenProxy() (string, error) {
	// 获取所有的代理软件
	proxySoftwares := DefaultSoftManager.GetProxySoftwares()

	// 尝试逐个启动，如果有启动成功，就返回
	for _, software := range proxySoftwares {
		soft, err := DefaultSoftManager.GetSoftware(software)
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
func (p *ProxyManager) CloseProxy(softwareName string) error {
	soft, err := DefaultSoftManager.GetSoftware(softwareName)
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
func (p *ProxyManager) checkEnvProxy() bool {
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
func (p *ProxyManager) checkSystemProxy() bool {
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
func (p *ProxyManager) checkHTTPClientProxy() bool {
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
