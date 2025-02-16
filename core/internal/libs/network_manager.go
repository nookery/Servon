package libs

import (
	"bufio"
	"context"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	netutil "github.com/shirou/gopsutil/v3/net"
)

type NetworkManager struct {
	NetworkStats *NetworkStats
}

func NewNetworkManager() *NetworkManager {
	return &NetworkManager{}
}

type NetworkStats struct {
	DownloadSpeed int64    `json:"download_speed"` // 下载速度（字节/秒）
	UploadSpeed   int64    `json:"upload_speed"`   // 上传速度（字节/秒）
	IPAddresses   IPConfig `json:"ip_addresses"`   // IP地址配置信息
}

// IPConfig 存储所有IP相关信息
type IPConfig struct {
	LocalIPs     []LocalIPInfo `json:"local_ips"`     // 本地IP信息
	PublicIP     string        `json:"public_ip"`     // 公网IP
	PublicIPv6   string        `json:"public_ipv6"`   // 公网IPv6
	DNSServers   []string      `json:"dns_servers"`   // DNS服务器
	NetworkCards []NetworkCard `json:"network_cards"` // 网卡信息
}

// LocalIPInfo 本地IP详细信息
type LocalIPInfo struct {
	IP        string `json:"ip"`        // IP地址
	Interface string `json:"interface"` // 网卡接口名称
	IsIPv6    bool   `json:"is_ipv6"`   // 是否是IPv6
	NetMask   string `json:"netmask"`   // 子网掩码
}

// NetworkCard 网卡信息
type NetworkCard struct {
	Name       string   `json:"name"`        // 网卡名称
	MACAddress string   `json:"mac_address"` // MAC地址
	IsUp       bool     `json:"is_up"`       // 是否启用
	MTU        int      `json:"mtu"`         // MTU值
	IPs        []string `json:"ips"`         // 分配的IP地址
}

var (
	lastStats     []netutil.IOCountersStat
	lastStatsTime time.Time
)

// GetNetworkResources 获取网络资源使用情况
func (p *NetworkManager) GetNetworkResources() (*NetworkStats, error) {
	// 获取当前网络统计信息
	currentStats, err := netutil.IOCounters(false) // false表示获取所有网卡的总和
	if err != nil {
		return nil, err
	}

	currentTime := time.Now()

	// 如果是第一次获取数据
	if lastStats == nil {
		lastStats = currentStats
		lastStatsTime = currentTime
		return &NetworkStats{
			DownloadSpeed: 0,
			UploadSpeed:   0,
		}, nil
	}

	// 计算时间差（秒）
	duration := currentTime.Sub(lastStatsTime).Seconds()

	// 计算速度
	var totalBytesRecv int64
	var totalBytesSent int64
	var lastBytesRecv int64
	var lastBytesSent int64

	for _, stat := range currentStats {
		totalBytesRecv += int64(stat.BytesRecv)
		totalBytesSent += int64(stat.BytesSent)
	}

	for _, stat := range lastStats {
		lastBytesRecv += int64(stat.BytesRecv)
		lastBytesSent += int64(stat.BytesSent)
	}

	// 计算每秒的速度
	downloadSpeed := int64(float64(totalBytesRecv-lastBytesRecv) / duration)
	uploadSpeed := int64(float64(totalBytesSent-lastBytesSent) / duration)

	// 更新上次的统计信息
	lastStats = currentStats
	lastStatsTime = currentTime

	return &NetworkStats{
		DownloadSpeed: downloadSpeed,
		UploadSpeed:   uploadSpeed,
	}, nil
}

// GetIPConfig 获取完整的IP配置信息
func (p *NetworkManager) GetIPConfig() (*IPConfig, error) {
	ipConfig := &IPConfig{}

	// 获取本地IP
	localIPs, err := p.getLocalIPs()
	if err != nil {
		return nil, err
	}
	ipConfig.LocalIPs = localIPs

	// 获取公网IP
	publicIP, err := p.getPublicIP()
	if err == nil { // 即使获取失败也不影响其他信息
		ipConfig.PublicIP = publicIP
	}

	// 获取公网IPv6
	publicIPv6, err := p.getPublicIPv6()
	if err == nil {
		ipConfig.PublicIPv6 = publicIPv6
	}

	// 获取DNS服务器
	dnsServers, err := p.getDNSServers()
	if err == nil {
		ipConfig.DNSServers = dnsServers
	}

	// 获取网卡信息
	networkCards, err := p.getNetworkCards()
	if err == nil {
		ipConfig.NetworkCards = networkCards
	}

	return ipConfig, nil
}

// getLocalIPs 获取本地IP地址
func (p *NetworkManager) getLocalIPs() ([]LocalIPInfo, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var localIPs []LocalIPInfo
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			var mask net.IPMask

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
				mask = v.Mask
			case *net.IPAddr:
				ip = v.IP
				mask = ip.DefaultMask()
			}

			if ip == nil {
				continue
			}

			// 排除回环地址
			if ip.IsLoopback() {
				continue
			}

			localIPs = append(localIPs, LocalIPInfo{
				IP:        ip.String(),
				Interface: iface.Name,
				IsIPv6:    ip.To4() == nil,
				NetMask:   net.IP(mask).String(),
			})
		}
	}

	return localIPs, nil
}

// getPublicIP 获取公网IPv4地址
func (p *NetworkManager) getPublicIP() (string, error) {
	// 使用多个IP查询服务，提高可靠性
	clients := []struct {
		URL     string
		Timeout time.Duration
	}{
		{"https://api.ipify.org", 5 * time.Second},
		{"https://ifconfig.me/ip", 5 * time.Second},
		{"https://api.ip.sb/ip", 5 * time.Second},
	}

	for _, client := range clients {
		ctx, cancel := context.WithTimeout(context.Background(), client.Timeout)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "GET", client.URL, nil)
		if err != nil {
			continue
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		ip, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		return strings.TrimSpace(string(ip)), nil
	}

	return "", errors.New("failed to get public IP")
}

// getPublicIPv6 获取公网IPv6地址
func (p *NetworkManager) getPublicIPv6() (string, error) {
	// 类似getPublicIP的实现，但使用支持IPv6的服务
	// clients := []struct {
	// 	URL     string
	// 	Timeout time.Duration
	// }{
	// 	{"https://api6.ipify.org", 5 * time.Second},
	// 	{"https://v6.ident.me", 5 * time.Second},
	// }

	// 实现类似getPublicIP的逻辑
	// ... 实现代码 ...

	return "", errors.New("failed to get public IPv6")
}

// getDNSServers 获取DNS服务器列表
func (p *NetworkManager) getDNSServers() ([]string, error) {
	// 在Linux系统下读取/etc/resolv.conf
	file, err := os.Open("/etc/resolv.conf")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var servers []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "nameserver") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				servers = append(servers, fields[1])
			}
		}
	}

	return servers, scanner.Err()
}

// getNetworkCards 获取网卡信息
func (p *NetworkManager) getNetworkCards() ([]NetworkCard, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var cards []NetworkCard
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		var ips []string
		for _, addr := range addrs {
			ips = append(ips, addr.String())
		}

		card := NetworkCard{
			Name:       iface.Name,
			MACAddress: iface.HardwareAddr.String(),
			IsUp:       iface.Flags&net.FlagUp != 0,
			MTU:        iface.MTU,
			IPs:        ips,
		}
		cards = append(cards, card)
	}

	return cards, nil
}
