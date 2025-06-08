// Package command_util 提供命令行工具和选项管理功能
//
// 这个组件封装了与命令行相关的工具函数和类型定义，
// 包括命令选项结构、参数处理等功能。
package command_util

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var DefaultCommandUtil = &CommandUtil{}
var JoinArgs = DefaultCommandUtil.JoinArgs
var NewCommand = DefaultCommandUtil.NewCommand

type CommandUtil struct{}
type CommandOptions struct {
	Use     string
	Short   string
	Args    cobra.PositionalArgs
	Run     func(cmd *cobra.Command, args []string)
	Aliases []string
}

// JoinArgs 将参数数组连接成字符串，去掉方括号
func (c *CommandUtil) JoinArgs(args []string) string {
	result := ""
	for i, arg := range args {
		if i > 0 {
			result += " "
		}
		result += arg
	}
	return result
}

// NewCommand 创建新的命令
func (c *CommandUtil) NewCommand(options CommandOptions) *cobra.Command {
	return &cobra.Command{
		Use:     options.Use,
		Short:   options.Short,
		Args:    options.Args,
		Run:     options.Run,
		Aliases: options.Aliases,
	}
}

// ExecuteCommand 执行系统命令
func (c *CommandUtil) ExecuteCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("command failed: %v, output: %s", err, string(output))
	}
	return string(output), nil
}

// ExecuteCommandInDir 在指定目录执行系统命令
func (c *CommandUtil) ExecuteCommandInDir(dir, command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("command failed: %v, output: %s", err, string(output))
	}
	return string(output), nil
}

// ExecuteCommandWithEnv 执行带环境变量的系统命令
func (c *CommandUtil) ExecuteCommandWithEnv(env []string, command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Env = append(os.Environ(), env...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("command failed: %v, output: %s", err, string(output))
	}
	return string(output), nil
}

// IsCommandAvailable 检查命令是否可用
func (c *CommandUtil) IsCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// SplitCommand 分割命令字符串为命令和参数
func (c *CommandUtil) SplitCommand(commandLine string) (string, []string) {
	parts := strings.Fields(commandLine)
	if len(parts) == 0 {
		return "", nil
	}
	return parts[0], parts[1:]
}

// QuoteArgs 为参数添加引号（如果包含空格）
func (c *CommandUtil) QuoteArgs(args []string) []string {
	quoted := make([]string, len(args))
	for i, arg := range args {
		if strings.Contains(arg, " ") {
			quoted[i] = fmt.Sprintf("\"%s\"", arg)
		} else {
			quoted[i] = arg
		}
	}
	return quoted
}

// BuildCommandString 构建完整的命令字符串
func (c *CommandUtil) BuildCommandString(command string, args []string) string {
	if len(args) == 0 {
		return command
	}
	return fmt.Sprintf("%s %s", command, c.JoinArgs(c.QuoteArgs(args)))
}

// ValidateCommand 验证命令是否有效
func (c *CommandUtil) ValidateCommand(command string) error {
	if command == "" {
		return fmt.Errorf("command cannot be empty")
	}
	if !c.IsCommandAvailable(command) {
		return fmt.Errorf("command '%s' not found", command)
	}
	return nil
}

// GetCommandPath 获取命令的完整路径
func (c *CommandUtil) GetCommandPath(command string) (string, error) {
	return exec.LookPath(command)
}

// ExecuteCommandAsync 异步执行命令
func (c *CommandUtil) ExecuteCommandAsync(command string, args ...string) (*exec.Cmd, error) {
	cmd := exec.Command(command, args...)
	err := cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start command: %v", err)
	}
	return cmd, nil
}

// WaitForCommand 等待异步命令完成
func (c *CommandUtil) WaitForCommand(cmd *exec.Cmd) error {
	return cmd.Wait()
}

// KillCommand 终止命令进程
func (c *CommandUtil) KillCommand(cmd *exec.Cmd) error {
	if cmd.Process != nil {
		return cmd.Process.Kill()
	}
	return fmt.Errorf("no process to kill")
}