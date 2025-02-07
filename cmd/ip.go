package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var IPCmd = &cobra.Command{
	Use:   "ip",
	Short: "Get current server IP address",
	Long:  "Get the current server's public and local IP addresses",
	RunE: func(cmd *cobra.Command, args []string) error {
		ips, err := getAllIPs()
		if err != nil {
			return fmt.Errorf("failed to get IP addresses: %v", err)
		}

		// 分类显示IP
		var publicIPs, privateIPs []string
		for _, ip := range ips {
			if isPrivateIP(ip) {
				privateIPs = append(privateIPs, ip)
			} else {
				publicIPs = append(publicIPs, ip)
			}
		}

		// 显示公网IP
		if len(publicIPs) > 0 {
			color.Green("Public IPs:")
			for _, ip := range publicIPs {
				fmt.Printf("  - %s\n", ip)
			}
		}

		// 显示内网IP
		if len(privateIPs) > 0 {
			color.Green("\nPrivate IPs:")
			for _, ip := range privateIPs {
				fmt.Printf("  - %s\n", ip)
			}
		}

		return nil
	},
}

func getAllIPs() ([]string, error) {
	var ips []string

	// 获取公网IP
	publicIP, err := getPublicIP()
	if err == nil && publicIP != "" {
		ips = append(ips, publicIP)
	}

	// 获取本地IP
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		// 跳过down的接口和loopback
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			// 只处理IP地址
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ips = append(ips, ipnet.IP.String())
				}
			}
		}
	}

	return ips, nil
}

func isPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	// RFC 1918 私有网络:
	// 10.0.0.0/8
	// 172.16.0.0/12
	// 192.168.0.0/16
	privateRanges := []struct {
		start net.IP
		end   net.IP
	}{
		{
			start: net.ParseIP("10.0.0.0"),
			end:   net.ParseIP("10.255.255.255"),
		},
		{
			start: net.ParseIP("172.16.0.0"),
			end:   net.ParseIP("172.31.255.255"),
		},
		{
			start: net.ParseIP("192.168.0.0"),
			end:   net.ParseIP("192.168.255.255"),
		},
	}

	for _, r := range privateRanges {
		if bytes.Compare(ip, r.start) >= 0 && bytes.Compare(ip, r.end) <= 0 {
			return true
		}
	}

	return false
}

// 新增获取公网IP的函数
func getPublicIP() (string, error) {
	// 使用多个IP查询服务以提高可靠性
	urls := []string{
		"https://api.ipify.org",
		"https://ifconfig.me/ip",
		"https://api.ip.sb/ip",
	}

	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		ip, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		// 清理返回的IP地址（去除空格和换行符）
		ipStr := strings.TrimSpace(string(ip))
		if net.ParseIP(ipStr) != nil {
			return ipStr, nil
		}
	}

	return "", fmt.Errorf("failed to get public IP")
}
