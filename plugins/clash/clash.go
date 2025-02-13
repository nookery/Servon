package clash

import (
	"fmt"
	"os"
	"os/exec"
	"servon/core"
	"servon/core/contract"
	"strings"
)

const repoUrl = "https://github.com/wnlen/clash-for-linux.git"

func Setup(core *core.Core) {
	plugin := NewClash(core)

	core.RegisterSoftware("clash", plugin)
}

type Clash struct {
	info      contract.SoftwareInfo
	targetDir string
	*core.Core
}

func NewClash(core *core.Core) contract.SuperSoft {
	return &Clash{
		info: contract.SoftwareInfo{
			Name:        "clash",
			Description: "A rule-based tunnel in Go",
		},
		targetDir: core.GetDataRootFolder() + "/clash-for-linux",
		Core:      core,
	}
}

// Install 安装 Clash，并发送日志到日志通道
func (c *Clash) Install(logChan chan<- string) error {
	osType := c.GetOSType()

	if logChan == nil {
		logChan = make(chan string, 100)
	}

	c.PrintInfof("ClashPlugin: 检测到操作系统: %s", osType)
	c.PrintInfof("ClashPlugin: 开始安装 Clash...")

	switch osType {
	case core.Ubuntu, core.Debian:
		c.PrintInfof("ClashPlugin: 清理目标文件夹 - %s", c.targetDir)
		err := os.RemoveAll(c.targetDir)
		if err != nil {
			return fmt.Errorf("清理目标文件夹失败: %s", err)
		}
		c.PrintInfof("ClashPlugin: 目标文件夹清理完成 - %s", c.targetDir)

		// Clone clash-for-linux repository
		c.PrintInfof("ClashPlugin: 克隆 clash-for-linux 仓库 -> %s", repoUrl)
		err = c.RunShell("git", "clone", repoUrl, c.targetDir)
		if err != nil {
			return fmt.Errorf("克隆仓库失败: %s", err)
		}

		// 确保在 RunShellAndSendLog 完成后发送成功消息
		c.PrintInfof("ClashPlugin: 克隆仓库成功")
	case core.CentOS, core.RedHat:
		errMsg := "暂不支持在 RHEL 系统上安装 Clash"
		c.PrintErrorMessage(errMsg)
		return fmt.Errorf("%s", errMsg)

	default:
		errMsg := fmt.Sprintf("不支持的操作系统: %s", osType)
		c.PrintErrorMessage(errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	// 确保在返回前发送完成消息
	c.PrintInfof("ClashPlugin: 安装完成")
	return nil
}

func (c *Clash) Uninstall(logChan chan<- string) error {
	outputChan := make(chan string, 100)

	go func() {
		defer close(outputChan)

		// 停止服务
		outputChan <- "停止 Clash 服务..."
		stopCmd := exec.Command("sudo", "systemctl", "stop", "clash")
		if output, err := stopCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("停止服务失败:\n%s", string(output))
		}

		// 禁用服务
		disableCmd := exec.Command("sudo", "systemctl", "disable", "clash")
		if output, err := disableCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("禁用服务失败:\n%s", string(output))
		}

		// 删除服务文件
		rmServiceCmd := exec.Command("sudo", "rm", "/etc/systemd/system/clash.service")
		if output, err := rmServiceCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("删除服务文件失败:\n%s", string(output))
		}

		// 删除二进制文件
		rmBinCmd := exec.Command("sudo", "rm", "/usr/local/bin/clash")
		if output, err := rmBinCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("删除二进制文件失败:\n%s", string(output))
		}

		// 删除配置目录
		rmConfigCmd := exec.Command("sudo", "rm", "-rf", "/etc/clash")
		if output, err := rmConfigCmd.CombinedOutput(); err != nil {
			outputChan <- fmt.Sprintf("删除配置目录失败:\n%s", string(output))
		}

		outputChan <- "Clash 卸载完成"
	}()

	return nil
}

func (c *Clash) GetStatus() (map[string]string, error) {
	// 检查二进制文件是否存在
	if _, err := os.Stat("/usr/local/bin/clash"); os.IsNotExist(err) {
		return map[string]string{
			"status":  "not_installed",
			"version": "",
		}, nil
	}

	status := "stopped"
	if c.IsActive("clash") {
		status = "running"
	}

	// 获取版本
	version := ""
	verCmd := exec.Command("clash", "-v")
	if verOutput, err := verCmd.CombinedOutput(); err == nil {
		version = strings.TrimSpace(string(verOutput))
	}

	return map[string]string{
		"status":  status,
		"version": version,
	}, nil
}

func (c *Clash) GetInfo() contract.SoftwareInfo {
	return c.info
}

func (c *Clash) Start(logChan chan<- string) error {
	// 检查env文件是否配置
	envFile := c.targetDir + "/.env"
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		return fmt.Errorf("env文件不存在，请先配置env文件")
	}

	c.PrintInfo("开始启动 Clash")
	c.PrintInfof("配置文件 %s", envFile)

	// 读取配置文件
	envContent, err := os.ReadFile(envFile)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %s", err)
	}

	// 检查 CLASH_URL 是否已配置
	if strings.Contains(string(envContent), "CLASH_URL='更改为你的clash订阅地址'") ||
		strings.Contains(string(envContent), "CLASH_URL=your_subscription_url_here") ||
		!strings.Contains(string(envContent), "CLASH_URL") {

		c.PrintInfo("请输入你的 Clash 订阅地址：")
		var subscriptionURL string
		fmt.Scanln(&subscriptionURL)

		if subscriptionURL == "" {
			return fmt.Errorf("订阅地址不能为空")
		}

		// 更新配置文件
		newContent := strings.Replace(string(envContent),
			"CLASH_URL='更改为你的clash订阅地址'",
			fmt.Sprintf("CLASH_URL='%s'", subscriptionURL), -1)
		newContent = strings.Replace(newContent,
			"CLASH_URL=your_subscription_url_here",
			fmt.Sprintf("CLASH_URL='%s'", subscriptionURL), -1)

		err = os.WriteFile(envFile, []byte(newContent), 0644)
		if err != nil {
			return fmt.Errorf("更新配置文件失败: %s", err)
		}

		c.PrintSuccess("订阅地址已更新")
	}

	// 启动服务
	err = c.RunShell("bash", c.targetDir+"/start.sh")
	if err != nil {
		return fmt.Errorf("启动失败: %s", err)
	}

	return nil
}

func (c *Clash) Stop() error {
	return fmt.Errorf("not-implemented")
}

func (c *Clash) Reload() error {
	return fmt.Errorf("not-implemented")
}
