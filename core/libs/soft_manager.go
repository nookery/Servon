package libs

import (
	"fmt"
	"servon/core/contract"

	"github.com/spf13/cobra"
)

type SoftManager struct {
	Softwares map[string]contract.SuperSoft
}

func NewSoftManager() *SoftManager {
	return &SoftManager{
		Softwares: make(map[string]contract.SuperSoft),
	}
}

// newInfoCmd 返回 info 子命令
func (p *SoftManager) newInfoCmd() *cobra.Command {
	return NewCommand(CommandOptions{
		Use:   "info",
		Short: "显示软件详细信息",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				printer.PrintError(fmt.Errorf("\n❌ 缺少软件名称参数"))
				fmt.Println("\n用法:")
				printer.PrintYellow("  servon software info ")
				fmt.Println("[软件名称]")

				// 显示支持的软件列表
				names := p.GetAllSoftware()
				printer.PrintList(names, "支持的软件列表")

				return
			}

			name := args[0]

			// 检查软件是否支持
			supported := false
			for _, sw := range p.GetAllSoftware() {
				if sw == name {
					supported = true
					break
				}
			}

			if !supported {
				printer.PrintErrorMessage(fmt.Sprintf("不支持的软件: %s", name))
				fmt.Println("\n支持的软件:")
				names := p.GetAllSoftware()
				for _, name := range names {
					printer.PrintList([]string{name}, "支持的软件")
				}
				return
			}

			// 获取软件状态
			status, err := p.GetSoftwareStatus(name)
			if err != nil {
				printer.PrintErrorMessage(fmt.Sprintf("获取软件状态失败: %v", err))
				return
			}

			// 显示软件信息
			fmt.Println()
			printer.PrintCyan("%s", fmt.Sprintf("📦 %s\n", name))
			fmt.Println()

			// 显示安装状态
			printer.PrintWhite("状态: ")
			switch status["status"] {
			case "running":
				printer.PrintGreen("运行中")
			case "stopped":
				printer.PrintYellow("已停止")
			case "not_installed":
				printer.PrintRed("未安装")
			default:
				printer.PrintWhite(status["status"])
			}

			// 显示版本信息
			if version := status["version"]; version != "" {
				printer.PrintWhite("版本: ")
				printer.PrintWhite(version)
			}

			fmt.Println()
		},
	})
}

// GetSoftwareCommand 返回 software 命令
func (p *SoftManager) GetSoftwareCommand() *cobra.Command {
	cmd := NewCommand(CommandOptions{
		Use:   "software",
		Short: "软件管理",
	})

	cmd.AddCommand(p.newListCmd())
	cmd.AddCommand(p.newInstallCmd())
	cmd.AddCommand(p.newInfoCmd())
	cmd.AddCommand(p.newStartCmd())
	cmd.AddCommand(p.newStopCmd())
	cmd.AddCommand(p.newUninstallCmd())

	return cmd
}

// newStartCmd 返回 start 子命令
func (p *SoftManager) newStartCmd() *cobra.Command {
	return NewCommand(CommandOptions{
		Use:   "start",
		Short: "启动指定的软件",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				printer.PrintErrorMessage("缺少软件名称参数")
				fmt.Println("\n用法:")
				printer.PrintYellow("  servon software start ")
				fmt.Println("[软件名称]")

				// 显示支持的软件列表
				names := p.GetAllSoftware()
				fmt.Println("\n支持的软件:")
				for _, name := range names {
					printer.PrintList([]string{name}, "支持的软件")
				}

				fmt.Println("\n示例:")
				printer.PrintCyan("  servon software start nginx")
				printer.PrintCyan("  servon software start mysql")
				return
			}

			name := args[0]

			// 检查软件是否支持
			supported := false
			for _, sw := range p.GetAllSoftware() {
				if sw == name {
					supported = true
					break
				}
			}

			if !supported {
				printer.PrintErrorMessage(fmt.Sprintf("不支持的软件: %s", name))
				fmt.Println("\n支持的软件:")
				names := p.GetAllSoftware()
				for _, name := range names {
					printer.PrintList([]string{name}, "支持的软件")
				}
				return
			}

			// 开始启动
			Infoln("🚀 %s 启动中 ...", name)

			err := p.StartSoftware(name, nil)
			if err != nil {
				Infoln("❌ %s 启动失败", name)
				Error("%s", err)
				return
			}

			Infoln("✅ %s 启动成功！", name)
		},
	})
}

func (p *SoftManager) newStopCmd() *cobra.Command {
	return NewCommand(CommandOptions{
		Use:   "stop",
		Short: "停止指定的软件",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				printer.PrintErrorMessage("缺少软件名称参数")
				fmt.Println("\n用法:")
				printer.PrintYellow("  servon software stop ")
				fmt.Println("[软件名称]")

				// 显示支持的软件列表
				names := p.GetAllSoftware()
				fmt.Println("\n支持的软件:")
				for _, name := range names {
					printer.PrintList([]string{name}, "支持的软件")
				}

				fmt.Println("\n示例:")
				printer.PrintCyan("  servon software stop caddy")
				printer.PrintCyan("  servon software stop clash")
				return
			}

			name := args[0]

			// 检查软件是否支持
			supported := false
			for _, sw := range p.GetAllSoftware() {
				if sw == name {
					supported = true
					break
				}
			}

			if !supported {
				printer.PrintErrorMessage(fmt.Sprintf("不支持的软件: %s", name))
				fmt.Println("\n支持的软件:")
				names := p.GetAllSoftware()
				for _, name := range names {
					printer.PrintList([]string{name}, "支持的软件")
				}
				return
			}

			// 开始停止
			Infoln("🛑 %s 停止中 ...", name)

			err := p.StopSoftware(name)
			if err != nil {
				Infoln("❌ %s 停止失败", name)
				Error("%s", err)
				return
			}

			Infoln("✅ %s 已停止！", name)
		},
	})
}

// newListCmd 返回 list 子命令
func (p *SoftManager) newListCmd() *cobra.Command {
	return NewCommand(CommandOptions{
		Use:     "list",
		Short:   "显示支持的软件列表",
		Aliases: []string{"l"},
		Run: func(cmd *cobra.Command, args []string) {
			names := p.GetAllSoftware()

			printer.PrintList(names, "支持的软件列表")
		},
	})
}

func (p *SoftManager) newInstallCmd() *cobra.Command {
	cmd := NewCommand(CommandOptions{
		Use:     "install",
		Short:   "安装指定的软件",
		Args:    cobra.ExactArgs(1),
		Aliases: []string{"i"},
		Run: func(cmd *cobra.Command, args []string) {
			p.Install(args[0], nil)
		},
	})

	return cmd
}

// newUninstallCmd 返回 uninstall 子命令
func (p *SoftManager) newUninstallCmd() *cobra.Command {
	return NewCommand(CommandOptions{
		Use:     "uninstall",
		Short:   "卸载指定的软件",
		Aliases: []string{"u", "remove"},
		Run: func(cmd *cobra.Command, args []string) {
			p.UninstallSoftware(args[0], nil)
		},
	})
}

// Install 安装软件, 如果提供了日志通道则输出日志
func (c *SoftManager) Install(name string, logChan chan<- string) error {
	software, ok := c.Softwares[name]
	if !ok {
		registeredSoftwares := make([]string, 0, len(c.Softwares))
		for name := range c.Softwares {
			registeredSoftwares = append(registeredSoftwares, name)
		}

		printer.PrintList(registeredSoftwares, "可用的软件")
		return printer.PrintAndReturnError(fmt.Sprintf("软件 %s 未注册", name))
	}
	return software.Install(logChan)
}

// UninstallSoftware 卸载软件
func (c *SoftManager) UninstallSoftware(name string, logChan chan<- string) error {
	software, ok := c.Softwares[name]
	if !ok {
		return printer.PrintAndReturnError(fmt.Sprintf("软件 %s 未注册", name))
	}
	return software.Uninstall(logChan)
}

// StartSoftware 启动软件
func (c *SoftManager) StartSoftware(name string, logChan chan<- string) error {
	software, ok := c.Softwares[name]
	if !ok {
		return printer.PrintAndReturnError(fmt.Sprintf("软件 %s 未注册", name))
	}
	return software.Start(logChan)
}

// StopSoftware 停止软件
func (c *SoftManager) StopSoftware(name string) error {
	software, ok := c.Softwares[name]
	if !ok {
		return printer.PrintAndReturnError(fmt.Sprintf("软件 %s 未注册", name))
	}
	return software.Stop()
}

// GetSoftwareStatus 获取软件状态
func (c *SoftManager) GetSoftwareStatus(name string) (map[string]string, error) {
	software, ok := c.Softwares[name]
	if !ok {
		return nil, printer.PrintAndReturnError(fmt.Sprintf("软件 %s 未注册", name))
	}
	return software.GetStatus()
}

// RegisterSoftware 注册软件
func (c *SoftManager) RegisterSoftware(name string, software contract.SuperSoft) error {
	if _, exists := c.Softwares[name]; exists {
		return printer.PrintAndReturnError(fmt.Sprintf("软件 %s 已注册", name))
	}
	c.Softwares[name] = software
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
