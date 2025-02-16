package commands

import (
	"servon/core/internal/libs"

	"strings"

	"github.com/spf13/cobra"
)

var NetworkManager = libs.DefaultNetworkManager

var IPCmd = &cobra.Command{
	Use:   "ip",
	Short: "Get IP-related information",
	Long:  `Display various IP-related information including local IPs, public IP, and network interfaces`,
}

func init() {
	// Add subcommands
	IPCmd.AddCommand(localCmd)
	IPCmd.AddCommand(publicCmd)
	IPCmd.AddCommand(interfacesCmd)
}

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Show local IP addresses",
	RunE: func(cmd *cobra.Command, args []string) error {
		ipConfig, err := NetworkManager.GetIPConfig()
		if err != nil {
			return err
		}

		for _, ip := range ipConfig.LocalIPs {
			cmd.Printf("Interface: %s\n", ip.Interface)
			cmd.Printf("IP: %s\n", ip.IP)
			cmd.Printf("Netmask: %s\n", ip.NetMask)
			cmd.Printf("IPv6: %v\n", ip.IsIPv6)
			cmd.Println("---")
		}
		return nil
	},
}

var publicCmd = &cobra.Command{
	Use:   "public",
	Short: "Show public IP addresses",
	RunE: func(cmd *cobra.Command, args []string) error {
		ipConfig, err := NetworkManager.GetIPConfig()
		if err != nil {
			return err
		}

		if ipConfig.PublicIP != "" {
			cmd.Printf("Public IPv4: %s\n", ipConfig.PublicIP)
		}
		if ipConfig.PublicIPv6 != "" {
			cmd.Printf("Public IPv6: %s\n", ipConfig.PublicIPv6)
		}
		return nil
	},
}

var interfacesCmd = &cobra.Command{
	Use:   "interfaces",
	Short: "Show network interfaces information",
	RunE: func(cmd *cobra.Command, args []string) error {
		ipConfig, err := NetworkManager.GetIPConfig()
		if err != nil {
			return err
		}

		for _, card := range ipConfig.NetworkCards {
			cmd.Printf("Name: %s\n", card.Name)
			cmd.Printf("MAC Address: %s\n", card.MACAddress)
			cmd.Printf("Status: %v\n", map[bool]string{true: "UP", false: "DOWN"}[card.IsUp])
			cmd.Printf("MTU: %d\n", card.MTU)
			cmd.Printf("IP Addresses: %s\n", strings.Join(card.IPs, ", "))
			cmd.Println("---")
		}
		return nil
	},
}
