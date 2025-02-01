package main

import (
	"encoding/json"
	"fmt"
	"os"
	"servon/internal/system"
	"servon/internal/web"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	port     int
	interval int
	format   string
)

// 定义颜色输出
var (
	titleColor    = color.New(color.FgHiCyan, color.Bold)
	successColor  = color.New(color.FgGreen)
	errorColor    = color.New(color.FgRed)
	warningColor  = color.New(color.FgYellow)
	infoColor     = color.New(color.FgHiBlue)
	headingColor  = color.New(color.FgHiMagenta)
	defaultColor  = color.New(color.FgWhite)
)

var rootCmd = &cobra.Command{
	Use:   "servon",
	Short: "Servon - A lightweight server management panel",
	Long: `Servon is a comprehensive server management panel that provides
both CLI and Web interface for managing your server, websites, and docker containers.`,
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Servon web server",
	Run: func(cmd *cobra.Command, args []string) {
		titleColor.Printf("Starting Servon web server on port %d...\n", port)
		server := web.NewServer(port)
		if err := server.Start(); err != nil {
			errorColor.Printf("Failed to start server: %v\n", err)
			os.Exit(1)
		}
	},
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display system information",
	Run: func(cmd *cobra.Command, args []string) {
		info, err := system.GetSystemInfo()
		if err != nil {
			errorColor.Printf("Failed to get system info: %v\n", err)
			os.Exit(1)
		}

		switch format {
		case "json":
			printJSON(info)
		case "plain":
			printPlain(info)
		default:
			printFormatted(info)
		}
	},
}

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor system resources in real-time",
	Run: func(cmd *cobra.Command, args []string) {
		titleColor.Printf("Monitoring system resources (interval: %ds, press Ctrl+C to stop)...\n\n", interval)
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()

		// 打印表头
		headingColor.Printf("%-20s %-10s %-15s %-15s\n", "Time", "CPU(%)", "Memory(%)", "Disk(%)")
		fmt.Println(strings.Repeat("-", 65))

		for {
			select {
			case <-ticker.C:
				info, err := system.GetSystemInfo()
				if err != nil {
					errorColor.Printf("Error: %v\n", err)
					continue
				}

				memPercent := float64(info.MemoryUsed) / float64(info.MemoryTotal) * 100
				diskPercent := float64(info.DiskUsed) / float64(info.DiskTotal) * 100
				
				// 根据使用率选择颜色
				cpuColor := getResourceColor(info.CPUUsage)
				memColor := getResourceColor(memPercent)
				diskColor := getResourceColor(diskPercent)

				defaultColor.Printf("%-20s ", time.Now().Format("15:04:05"))
				cpuColor.Printf("%-10.2f ", info.CPUUsage)
				memColor.Printf("%-15.2f ", memPercent)
				diskColor.Printf("%-15.2f\n", diskPercent)
			}
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		titleColor.Println("Servon v0.1.0")
	},
}

func init() {
	// 全局标志
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "formatted", "Output format (formatted|json|plain)")

	// serve 命令的标志
	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the web server on")

	// monitor 命令的标志
	monitorCmd.Flags().IntVarP(&interval, "interval", "i", 5, "Monitoring interval in seconds")

	// 添加命令
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(monitorCmd)
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		errorColor.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// 辅助函数：以 JSON 格式打印
func printJSON(info *system.SystemInfo) {
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		errorColor.Printf("Failed to marshal JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}

// 辅助函数：以纯文本格式打印
func printPlain(info *system.SystemInfo) {
	fmt.Printf("Hostname=%s\n", info.Hostname)
	fmt.Printf("Platform=%s\n", info.Platform)
	fmt.Printf("CPUUsage=%.2f\n", info.CPUUsage)
	fmt.Printf("MemoryTotal=%d\n", info.MemoryTotal)
	fmt.Printf("MemoryUsed=%d\n", info.MemoryUsed)
	fmt.Printf("DiskTotal=%d\n", info.DiskTotal)
	fmt.Printf("DiskUsed=%d\n", info.DiskUsed)
	fmt.Printf("Uptime=%d\n", info.Uptime)
}

// 辅助函数：以格式化文本打印
func printFormatted(info *system.SystemInfo) {
	memPercent := float64(info.MemoryUsed) / float64(info.MemoryTotal) * 100
	diskPercent := float64(info.DiskUsed) / float64(info.DiskTotal) * 100

	titleColor.Println("System Information:")
	titleColor.Println("==================")
	
	infoColor.Printf("%-12s", "Hostname:")
	defaultColor.Printf(" %s\n", info.Hostname)
	
	infoColor.Printf("%-12s", "Platform:")
	defaultColor.Printf(" %s\n", info.Platform)
	
	infoColor.Printf("%-12s", "CPU Usage:")
	getResourceColor(info.CPUUsage).Printf(" %.2f%%\n", info.CPUUsage)
	
	infoColor.Printf("%-12s", "Memory:")
	getResourceColor(memPercent).Printf(" %.2f%% ", memPercent)
	defaultColor.Printf("(Used: %s, Total: %s)\n", 
		formatBytes(info.MemoryUsed),
		formatBytes(info.MemoryTotal),
	)
	
	infoColor.Printf("%-12s", "Disk:")
	getResourceColor(diskPercent).Printf(" %.2f%% ", diskPercent)
	defaultColor.Printf("(Used: %s, Total: %s)\n",
		formatBytes(info.DiskUsed),
		formatBytes(info.DiskTotal),
	)
	
	infoColor.Printf("%-12s", "Uptime:")
	defaultColor.Printf(" %s\n", formatUptime(info.Uptime))
}

// 辅助函数：格式化字节大小
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// 辅助函数：格式化运行时间
func formatUptime(uptime uint64) string {
	days := uptime / (24 * 3600)
	hours := (uptime % (24 * 3600)) / 3600
	minutes := (uptime % 3600) / 60
	seconds := uptime % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

// 辅助函数：根据资源使用率返回对应的颜色
func getResourceColor(usage float64) *color.Color {
	switch {
	case usage >= 90:
		return color.New(color.FgRed, color.Bold)    // 危险
	case usage >= 70:
		return color.New(color.FgYellow)             // 警告
	case usage >= 50:
		return color.New(color.FgHiYellow)           // 注意
	default:
		return color.New(color.FgGreen)              // 正常
	}
} 