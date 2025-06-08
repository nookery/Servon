package user

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// RunShell 执行shell命令，返回错误和命令的组合输出（stdout和stderr）
func RunShell(command string, args ...string) (error, string) {
	cmd := exec.Command(command, args...)

	// 获取命令的组合输出
	output, err := cmd.CombinedOutput()

	return err, string(output)
}

// RunShellWithOutput 执行shell命令并返回标准输出
// 如果命令执行失败，error中会包含stderr的内容
func RunShellWithOutput(command string, args ...string) (error, string) {
	cmd := exec.Command(command, args...)

	// 分别捕获stdout和stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// 如果有错误，将stderr的内容包含在错误信息中
		return fmt.Errorf("%v: %s", err, stderr.String()), stdout.String()
	}

	// 返回stdout的内容
	return nil, strings.TrimSpace(stdout.String())
}
