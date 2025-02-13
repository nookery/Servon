package libs

import (
	"fmt"
	"servon/core/contract"

	"github.com/spf13/cobra"
)

type SoftManager struct {
	Softwares map[string]contract.SuperSoft
}

func newSoftManager() *SoftManager {
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
				DefaultPrinter.PrintError(fmt.Errorf("\n❌ 缺少软件名称参数"))
				fmt.Println("\n用法:")
				DefaultPrinter.PrintYellow("  servon software info ")
				fmt.Println("[软件名称]")

				// 显示支持的软件列表
				names := p.GetAllSoftware()
				DefaultPrinter.PrintList(names, "支持的软件列表")

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
				DefaultPrinter.PrintErrorMessage(fmt.Sprintf("不支持的软件: %s", name))
				fmt.Println("\n支持的软件:")
				names := p.GetAllSoftware()
				for _, name := range names {
					DefaultPrinter.PrintList([]string{name}, "支持的软件")
				}
				return
			}

			// 获取软件状态
			status, err := p.GetSoftwareStatus(name)
			if err != nil {
				DefaultPrinter.PrintErrorMessage(fmt.Sprintf("获取软件状态失败: %v", err))
				return
			}

			// 显示软件信息
			fmt.Println()
			DefaultPrinter.PrintCyan("%s", fmt.Sprintf("📦 %s\n", name))
			fmt.Println()

			// 显示安装状态
			DefaultPrinter.PrintWhite("状态: ")
			switch status["status"] {
			case "running":
				DefaultPrinter.PrintGreen("运行中")
			case "stopped":
				DefaultPrinter.PrintYellow("已停止")
			case "not_installed":
				DefaultPrinter.PrintRed("未安装")
			default:
				DefaultPrinter.PrintWhite("%s", status["status"])
			}

			// 显示版本信息
			if version := status["version"]; version != "" {
				DefaultPrinter.PrintWhite("版本: ")
				DefaultPrinter.PrintWhite(version)
			}

			fmt.Println()
		},
	})
}

// GetSoftwareCommand 返回 software 命令
func (p *SoftManager) GetSoftwareCommand() *cobra.Command {
	cmd := NewCommand(CommandOptions{
		Use:     "software",
		Short:   "软件管理",
		Aliases: []string{"soft"},
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
				DefaultPrinter.PrintErrorMessage("缺少软件名称参数")
				fmt.Println("\n用法:")
				DefaultPrinter.PrintYellow("  servon software start ")
				fmt.Println("[软件名称]")

				// 显示支持的软件列表
				names := p.GetAllSoftware()
				fmt.Println("\n支持的软件:")
				for _, name := range names {
					DefaultPrinter.PrintList([]string{name}, "支持的软件")
				}

				fmt.Println("\n示例:")
				DefaultPrinter.PrintCyan("  servon software start nginx")
				DefaultPrinter.PrintCyan("  servon software start mysql")
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
				DefaultPrinter.PrintErrorMessage(fmt.Sprintf("不支持的软件: %s", name))
				fmt.Println("\n支持的软件:")
				names := p.GetAllSoftware()
				for _, name := range names {
					DefaultPrinter.PrintList([]string{name}, "支持的软件")
				}
				return
			}

			// 开始启动
			DefaultPrinter.PrintInfo(fmt.Sprintf("🚀 %s 启动中 ...", name))

			err := p.StartSoftware(name)
			if err != nil {
				DefaultPrinter.PrintErrorf("❌ %s 启动失败", name)
				return
			}

			DefaultPrinter.PrintInfo(fmt.Sprintf("✅ %s 启动成功！", name))
		},
	})
}

func (p *SoftManager) newStopCmd() *cobra.Command {
	return NewCommand(CommandOptions{
		Use:   "stop",
		Short: "停止指定的软件",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				DefaultPrinter.PrintErrorMessage("缺少软件名称参数")
				fmt.Println("\n用法:")
				DefaultPrinter.PrintYellow("  servon software stop ")
				fmt.Println("[软件名称]")

				// 显示支持的软件列表
				names := p.GetAllSoftware()
				fmt.Println("\n支持的软件:")
				for _, name := range names {
					DefaultPrinter.PrintList([]string{name}, "支持的软件")
				}

				fmt.Println("\n示例:")
				DefaultPrinter.PrintCyan("  servon software stop caddy")
				DefaultPrinter.PrintCyan("  servon software stop clash")
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
				DefaultPrinter.PrintErrorMessage(fmt.Sprintf("不支持的软件: %s", name))
				fmt.Println("\n支持的软件:")
				names := p.GetAllSoftware()
				for _, name := range names {
					DefaultPrinter.PrintList([]string{name}, "支持的软件")
				}
				return
			}

			// 开始停止
			DefaultPrinter.PrintInfo(fmt.Sprintf("🛑 %s 停止中 ...", name))

			err := p.StopSoftware(name)
			if err != nil {
				DefaultPrinter.PrintErrorf("❌ %s 停止失败", name)
				DefaultPrinter.PrintError(err)
				return
			}

			DefaultPrinter.PrintInfo(fmt.Sprintf("✅ %s 已停止！", name))
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

			DefaultPrinter.PrintList(names, "支持的软件列表")
		},
	})
}

// newInstallCmd 返回 install 子命令
func (p *SoftManager) newInstallCmd() *cobra.Command {
	cmd := NewCommand(CommandOptions{
		Use:     "install",
		Short:   "安装指定的软件",
		Args:    cobra.ExactArgs(1),
		Aliases: []string{"i"},
		Run: func(cmd *cobra.Command, args []string) {
			if err := p.Install(args[0]); err != nil {
				DefaultPrinter.PrintErrorf("安装失败: %v", err)
			}
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
			p.UninstallSoftware(args[0])
		},
	})
}

// Install 安装软件, 如果提供了日志通道则输出日志
func (c *SoftManager) Install(name string) error {
	software, ok := c.Softwares[name]
	if !ok {
		registeredSoftwares := make([]string, 0, len(c.Softwares))
		for name := range c.Softwares {
			registeredSoftwares = append(registeredSoftwares, name)
		}

		DefaultPrinter.PrintList(registeredSoftwares, "可用的软件")
		return DefaultPrinter.PrintAndReturnError(fmt.Sprintf("软件 %s 未注册", name))
	}

	return software.Install()
}

// UninstallSoftware 卸载软件
func (c *SoftManager) UninstallSoftware(name string) error {
	software, ok := c.Softwares[name]
	if !ok {
		return DefaultPrinter.PrintAndReturnError(fmt.Sprintf("软件 %s 未注册", name))
	}
	return software.Uninstall()
}

// StartSoftware 启动软件
func (c *SoftManager) StartSoftware(name string) error {
	software, ok := c.Softwares[name]
	if !ok {
		return DefaultPrinter.PrintAndReturnError(fmt.Sprintf("软件 %s 未注册", name))
	}
	return software.Start()
}

// StopSoftware 停止软件
func (c *SoftManager) StopSoftware(name string) error {
	software, ok := c.Softwares[name]
	if !ok {
		return DefaultPrinter.PrintAndReturnError(fmt.Sprintf("软件 %s 未注册", name))
	}
	return software.Stop()
}

// GetSoftwareStatus 获取软件状态
func (c *SoftManager) GetSoftwareStatus(name string) (map[string]string, error) {
	software, ok := c.Softwares[name]
	if !ok {
		return nil, DefaultPrinter.PrintAndReturnError(fmt.Sprintf("软件 %s 未注册", name))
	}
	return software.GetStatus()
}

// RegisterSoftware 注册软件
func (c *SoftManager) RegisterSoftware(name string, software contract.SuperSoft) error {
	if _, exists := c.Softwares[name]; exists {
		return DefaultPrinter.PrintAndReturnError(fmt.Sprintf("软件 %s 已注册", name))
	}
	c.Softwares[name] = software
	return nil
}

// GetAllSoftware 获取所有软件
func (c *SoftManager) GetAllSoftware() []string {
	DefaultPrinter.PrintInfo("获取所有软件...")
	softwareNames := make([]string, 0, len(c.Softwares))
	for name := range c.Softwares {
		softwareNames = append(softwareNames, name)
	}
	return softwareNames
}
