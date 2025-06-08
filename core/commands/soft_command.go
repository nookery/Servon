package commands

import (
	"fmt"
	"servon/core/managers"

	"github.com/spf13/cobra"
)

// GetSoftwareCommand 返回 software 命令
func GetSoftwareCommand(p *managers.SoftManager) *cobra.Command {
	cmd := NewCommand(CommandOptions{
		Use:     "software",
		Short:   "软件管理",
		Aliases: []string{"soft"},
	})

	cmd.AddCommand(newListCmd(p))
	cmd.AddCommand(newInstallCmd(p))
	cmd.AddCommand(newInfoCmd(p))
	cmd.AddCommand(newStartCmd(p))
	cmd.AddCommand(newStopCmd(p))
	cmd.AddCommand(newUninstallCmd(p))

	return cmd
}

// newInfoCmd 返回 info 子命令
func newInfoCmd(p *managers.SoftManager) *cobra.Command {
	return NewCommand(CommandOptions{
		Use:   "info",
		Short: "显示软件详细信息",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				logger.Error(fmt.Errorf("\n❌ 缺少软件名称参数"))
				fmt.Println("\n用法:")
				logger.Info("  servon software info ")
				fmt.Println("[软件名称]")

				// 显示支持的软件列表
				names := p.GetAllSoftware()
				logger.ListWithTitle("支持的软件列表", names)
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
				logger.ErrorMessage(fmt.Sprintf("不支持的软件: %s", name))
				fmt.Println("\n支持的软件:")
				names := p.GetAllSoftware()
				for _, name := range names {
					logger.ListWithTitle("支持的软件", []string{name})
				}
				return
			}

			// 获取软件状态
			status, err := p.GetSoftwareStatus(name)
			if err != nil {
				logger.ErrorMessage(fmt.Sprintf("获取软件状态失败: %v", err))
				return
			}

			// 显示软件信息
			fmt.Println()
			PrintTitle(name)
			fmt.Println()

			// 对状态进行本地化处理
			if statusValue, exists := status["status"]; exists {
				statusText := map[string]string{
					"not_installed": "未安装",
					"installed":     "已安装",
					"running":       "运行中",
					"stopped":       "已停止",
					"error":         "异常",
				}
				if localText, ok := statusText[statusValue]; ok {
					status["status"] = localText
				}
			}

			// 使用 PrintKeyValues 输出所有状态信息
			logger.PrintKeyValues(status)
			fmt.Println()
		},
	})
}

// newUninstallCmd 返回 uninstall 子命令
func newUninstallCmd(p *managers.SoftManager) *cobra.Command {
	return NewCommand(CommandOptions{
		Use:     "uninstall",
		Short:   "卸载指定的软件",
		Aliases: []string{"u", "remove"},
		Run: func(cmd *cobra.Command, args []string) {
			p.UninstallSoftware(args[0])
		},
	})
}

// newStartCmd 返回 start 子命令
func newStartCmd(p *managers.SoftManager) *cobra.Command {
	return NewCommand(CommandOptions{
		Use:   "start",
		Short: "启动指定的软件",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				logger.ErrorMessage("缺少软件名称参数")
				fmt.Println("\n用法:")
				logger.Info("  servon software start ")
				fmt.Println("[软件名称]")

				// 显示支持的软件列表
				names := p.GetAllSoftware()
				fmt.Println("\n支持的软件:")
				for _, name := range names {
					logger.ListWithTitle("支持的软件", []string{name})
				}

				fmt.Println("\n示例:")
				logger.Info("  servon software start nginx")
				logger.Info("  servon software start mysql")
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
				logger.ErrorMessage(fmt.Sprintf("不支持的软件: %s", name))
				fmt.Println("\n支持的软件:")
				names := p.GetAllSoftware()
				for _, name := range names {
					logger.ListWithTitle("支持的软件", []string{name})
				}
				return
			}

			// 开始启动
			logger.Info(fmt.Sprintf("🚀 %s 启动中 ...", name))

			err := p.StartSoftware(name)
			if err != nil {
				logger.Errorf("%s 启动失败", name)
				PrintError(err)
				return
			}

			logger.Successf("✅ %s 启动成功！", name)
		},
	})
}

func newStopCmd(p *managers.SoftManager) *cobra.Command {
	return NewCommand(CommandOptions{
		Use:   "stop",
		Short: "停止指定的软件",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				logger.ErrorMessage("缺少软件名称参数")
				fmt.Println("\n用法:")
				logger.Info("  servon software stop ")
				fmt.Println("[软件名称]")

				// 显示支持的软件列表
				names := p.GetAllSoftware()
				fmt.Println("\n支持的软件:")
				for _, name := range names {
					logger.ListWithTitle("支持的软件", []string{name})
				}

				fmt.Println("\n示例:")
				logger.Info("  servon software stop caddy")
				logger.Info("  servon software stop clash")
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
				logger.ErrorMessage(fmt.Sprintf("不支持的软件: %s", name))
				fmt.Println("\n支持的软件:")
				names := p.GetAllSoftware()
				for _, name := range names {
					logger.ListWithTitle("支持的软件", []string{name})
				}
				return
			}

			// 开始停止
			PrintInfof("%s 停止中 ...", name)

			err := p.StopSoftware(name)
			if err != nil {
				PrintErrorf("%s 停止失败", name)
				PrintError(err)
				return
			}

			PrintSuccessf("%s 已停止！", name)
		},
	})
}

// newListCmd 返回 list 子命令
func newListCmd(p *managers.SoftManager) *cobra.Command {
	return NewCommand(CommandOptions{
		Use:     "list",
		Short:   "显示支持的软件列表",
		Aliases: []string{"l"},
		Run: func(cmd *cobra.Command, args []string) {
			names := p.GetAllSoftware()

			logger.ListWithTitle("支持的软件列表", names)
		},
	})
}

// newInstallCmd 返回 install 子命令
func newInstallCmd(p *managers.SoftManager) *cobra.Command {
	cmd := NewCommand(CommandOptions{
		Use:     "install",
		Short:   "安装指定的软件",
		Args:    cobra.ExactArgs(1),
		Aliases: []string{"i"},
		Run: func(cmd *cobra.Command, args []string) {
			if err := p.Install(args[0]); err != nil {
				logger.Errorf("安装失败: %v", err)
			}
		},
	})

	return cmd
}
