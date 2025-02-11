package libs

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"servon/core/templates"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var printer = DefaultPrinter

type CommandOptions struct {
	Use     string
	Short   string
	Args    cobra.PositionalArgs
	Run     func(cmd *cobra.Command, args []string)
	Aliases []string
}

type CommandManager struct {
	CommandOptions
	rootCmd *cobra.Command
}

func NewCommandManager() *CommandManager {
	commandManager := &CommandManager{}
	commandManager.rootCmd = commandManager.NewCommand(CommandOptions{
		Use:   "servon",
		Short: "Servon 是一个用于管理服务器的命令行工具",
	})
	return commandManager
}

// AddCommand 添加一个命令
func (c *CommandManager) AddCommand(cmd *cobra.Command) {
	c.rootCmd.AddCommand(cmd)
}

// CheckCommandArgs 检查命令参数
func (c *CommandManager) CheckCommandArgs(cmd *cobra.Command, args []string) error {
	// 如果命令没有设置 Args 要求，则至少需要一个参数
	if cmd.Args == nil {
		if len(args) == 0 {
			return fmt.Errorf("至少需要一个参数")
		}
		return nil
	}

	// 使用命令自带的参数验证
	err := cmd.Args(cmd, args)
	if err != nil {
		return fmt.Errorf("参数验证失败: %v", err)
	}

	printer.PrintInfo("参数验证成功")
	return nil
}

func (c *CommandManager) Execute(command string, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("command is required")
	}

	// 使用青色（Cyan）输出命令和参数，用空格连接参数
	color.Cyan("📺 %s %s", command, joinArgs(args))

	execCmd := exec.Command(command, args...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin

	return execCmd.Run()
}

func (c *CommandManager) ExecuteWithOutput(command string, args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("command is required")
	}

	execCmd := exec.Command(command, args...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin

	output, err := execCmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// joinArgs 将参数数组连接成字符串，去掉方括号
func joinArgs(args []string) string {
	result := ""
	for i, arg := range args {
		if i > 0 {
			result += " "
		}
		result += arg
	}
	return result
}

// NewCommand 创建一个标准化的命令
func (c *CommandManager) NewCommand(opts CommandOptions) *cobra.Command {
	setCustomErrPrefix := true
	setCustomUsageTemplate := true
	setCustomHelpFunc := true

	cmd := &cobra.Command{
		Use:           opts.Use,
		Short:         opts.Short,
		SilenceErrors: false,
		SilenceUsage:  false,
		Args:          opts.Args,
		Aliases:       opts.Aliases,
		PreRun: func(cmd *cobra.Command, args []string) {
			// libs.Infoln("🚀 开始执行命令 PreRun")
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// libs.Infoln("🚀 开始执行命令 PreRunE")
			return nil
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// libs.Infoln("🚀 开始执行命令 PersistentPreRun")
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// libs.Infoln("🚀 开始执行命令 PersistentPreRunE")
			return nil
		},
		Run: opts.Run,
		PostRun: func(cmd *cobra.Command, args []string) {
			// libs.Infoln("🎉 命令执行成功 PostRun")
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			// libs.Infoln("🎉 命令执行完成 PostRunE")
			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			// libs.Infoln("🎉 命令执行完成 PersistentPostRun")
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			// libs.Infoln("🎉 命令执行完成 PersistentPostRunE")
			return nil
		},
	}

	if setCustomErrPrefix {
		cmd.SetErrPrefix("❌ 发生了错误")
	}

	// 自定义错误处理
	cmd.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
		c.Printf("%s\n", color.New(color.FgHiRed).Sprintf("%s", "❌ 错误："+fmt.Sprintf("%v", err)))
		return nil
	})

	// 自定义帮助
	if setCustomHelpFunc {
		cmd.SetHelpFunc(func(c *cobra.Command, args []string) {
			c.Printf("%s\n", color.New(color.FgHiCyan).Sprintf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
			c.Printf("📌 命令: %s\n", color.New(color.FgHiYellow).Sprintf(c.Use))
			c.Printf("📝 描述: %s\n", color.New(color.FgHiGreen).Sprintf(c.Short))
			c.Printf("\n%s\n", color.New(color.FgHiBlue).Sprintf("🎯 参数列表:"))
			c.Printf("%s\n", color.New(color.FgHiCyan).Sprintf("%s", c.LocalFlags().FlagUsages()))
			c.Printf("%s\n", color.New(color.FgHiCyan).Sprintf("🎯 可用命令:"))
			for _, command := range c.Commands() {
				alias := ""
				if len(command.Aliases) > 0 {
					alias = "(" + joinArgs(command.Aliases) + ")"
				}

				nameAndAlias := ""
				if alias != "" {
					nameAndAlias = fmt.Sprintf("%s %s", command.Use, alias)
				} else {
					nameAndAlias = command.Use
				}
				c.Printf("  %-35s%s\n", color.New(color.FgHiCyan).Sprintf("%s", nameAndAlias), command.Short)
			}
			c.Printf("%s\n", color.New(color.FgHiCyan).Sprintf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
		})
	}

	// 自定义使用说明模板
	if setCustomUsageTemplate {
		cmd.SetUsageTemplate(templates.UsageTemplate())
	}

	return cmd
}

// StreamCommand 执行命令并打印输出
func (c *CommandManager) StreamCommand(cmd *exec.Cmd) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("获取标准输出失败: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("获取标准错误输出失败: %v", err)
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	return cmd.Run()
}

// GetRootCommand 获取根命令
func (c *CommandManager) GetRootCommand() *cobra.Command {
	return c.rootCmd
}
