package clash

import (
	"fmt"
	"os"
	"os/exec"
	"servon/core"
	"strings"
)

const repoUrl = "https://github.com/wnlen/clash-for-linux.git"

var softInfo = core.SoftwareInfo{
	Name:            "clash",
	Description:     "A rule-based tunnel in Go",
	IsProxySoftware: true,
}

func Setup(app *core.App) {
	plugin := NewClash(app)

	app.RegisterSoftware("clash", plugin)
}

type Clash struct {
	targetDir string
	*core.App
}

func NewClash(app *core.App) core.SuperSoft {
	return &Clash{
		targetDir: app.GetSoftwareRootFolder("clash"),
		App:       app,
	}
}

// Install 安装 Clash，并发送日志到日志通道
func (c *Clash) Install() error {
	osType := c.GetOSType()

	c.PrintInfof("安装 Clash，检测到操作系统: %s", osType)

	switch osType {
	case core.Ubuntu, core.Debian:
		c.PrintInfof("清理目标文件夹 - %s", c.targetDir)
		err := os.RemoveAll(c.targetDir)
		if err != nil {
			return fmt.Errorf("清理目标文件夹失败: %s", err)
		}
		c.PrintInfof("目标文件夹清理完成 - %s", c.targetDir)

		// 使用 go-git 克隆仓库
		c.PrintInfof("克隆 clash-for-linux 仓库 -> %s", repoUrl)
		err = c.GitClone(repoUrl, "master", c.targetDir)
		if err != nil {
			c.PrintErrorMessage(err.Error())
			return fmt.Errorf("克隆仓库失败: %s", err)
		}

		c.PrintSuccess("克隆仓库成功")
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
	c.PrintInfof("Clash 安装完成")
	return nil
}

// Uninstall 卸载 Clash
func (c *Clash) Uninstall() error {
	// 停止服务
	err := c.Stop()
	if err != nil {
		return fmt.Errorf("停止服务失败: %s", err)
	}

	// 删除服务文件
	rmServiceCmd := exec.Command("sudo", "rm", "/etc/systemd/system/clash.service")
	if err := rmServiceCmd.Run(); err != nil {
		return fmt.Errorf("删除服务文件失败: %s", err)
	}

	// 删除二进制文件
	rmBinCmd := exec.Command("sudo", "rm", "/usr/local/bin/clash")
	if err := rmBinCmd.Run(); err != nil {
		return fmt.Errorf("删除二进制文件失败: %s", err)
	}

	// 删除配置目录
	rmConfigCmd := exec.Command("sudo", "rm", "-rf", "/etc/clash")
	if err := rmConfigCmd.Run(); err != nil {
		return fmt.Errorf("删除配置目录失败: %s", err)
	}

	c.PrintSuccess("Clash 卸载完成")
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

func (c *Clash) GetInfo() core.SoftwareInfo {
	return softInfo
}

func (c *Clash) Start() error {
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
