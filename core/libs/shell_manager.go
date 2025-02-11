package libs

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
)

type ShellManager struct {
}

func NewShellManager() *ShellManager {
	return &ShellManager{}
}

func (s *ShellManager) RunShell(command string, args ...string) error {
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
