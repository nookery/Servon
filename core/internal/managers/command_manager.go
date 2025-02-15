package managers

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"servon/core/internal/utils"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var printer = DefaultPrinter

type CommandManager struct {
	utils.CommandOptions
	rootCmd *cobra.Command
}

func NewCommandManager(rootCmd *cobra.Command) *CommandManager {
	commandManager := &CommandManager{}
	commandManager.rootCmd = rootCmd
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
	color.Cyan("📺 %s %s", command, utils.JoinArgs(args))

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
