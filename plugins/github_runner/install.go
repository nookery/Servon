package github_runner

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"servon/core"
	"strings"
	"time"
)

// Install 安装 GitHub Runner
func (g *GitHubRunner) Install() error {
	// 先检查运行状态
	status, err := g.GetStatus()
	if err != nil {
		return fmt.Errorf("检查运行状态失败: %s", err)
	}
	if status["status"] == "running" {
		return fmt.Errorf("GitHub Runner 正在运行中，请先停止后再安装")
	}

	osType := g.GetOSType()
	g.SoftwareLogger.Infof("检测到操作系统: %s", osType)

	if g.IsProxyOn() {
		g.SoftwareLogger.Infof("检测到代理")
	}

	// 创建 runner 专用用户
	g.SoftwareLogger.Infof("创建 GitHub Runner 专用用户...")
	runnerUser := "github-runner"
	runnerPassword := "runner" + fmt.Sprint(time.Now().Unix()) // 生成随机密码

	exists, err := g.UserExists(runnerUser)
	if err != nil {
		return fmt.Errorf("检查用户失败: %s", err)
	}

	if !exists {
		if err := g.CreateUser(runnerUser, runnerPassword); err != nil {
			return fmt.Errorf("创建用户失败: %s", err)
		}
		g.SoftwareLogger.Successf("已创建专用用户: %s", runnerUser)
	}

	switch osType {
	case core.Ubuntu, core.Debian, core.CentOS, core.RedHat:
		g.SoftwareLogger.Infof("开始安装 GitHub Runner...")

		// 获取最新版本号
		resp, err := http.Get("https://api.github.com/repos/actions/runner/releases/latest")
		if err != nil {
			return fmt.Errorf("获取最新版本失败: %s", err)
		}
		defer resp.Body.Close()

		// 检查状态码
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			if resp.StatusCode == http.StatusForbidden && strings.Contains(string(body), "API rate limit exceeded") {
				return fmt.Errorf("GitHub API 调用次数超限，请稍后再试或使用 API 令牌。状态码: %d", resp.StatusCode)
			}
			return fmt.Errorf("获取版本信息失败，GitHub API 返回状态码: %d，响应内容: %s", resp.StatusCode, string(body))
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("读取响应失败: %s", err)
		}

		// 解析版本号
		var release struct {
			TagName string `json:"tag_name"`
		}
		if err := json.Unmarshal(body, &release); err != nil {
			return fmt.Errorf("解析版本信息失败: %s", err)
		}

		if release.TagName == "" {
			return fmt.Errorf("获取版本号失败：API 返回的版本号为空")
		}

		version := strings.TrimPrefix(release.TagName, "v")
		if version == "" {
			return fmt.Errorf("获取版本号失败：无效的版本号格式")
		}

		g.SoftwareLogger.Infof("最新版本: %s", version)

		// 清理本地文件夹
		g.SoftwareLogger.Infof("清理本地文件夹...")
		err = os.RemoveAll(g.targetDir)
		if err != nil {
			return fmt.Errorf("清理本地文件夹失败: %s", err)
		}

		// 创建目标目录
		g.SoftwareLogger.Infof("创建目标目录...")
		err = os.MkdirAll(g.targetDir, 0755)
		if err != nil {
			return fmt.Errorf("创建目录失败: %s", err)
		}

		// 获取系统架构
		err, arch := g.RunShellWithOutput("uname", "-m")
		if err != nil {
			return fmt.Errorf("获取系统架构失败: %s", err)
		}
		arch = strings.TrimSpace(arch)

		// 确定下载的架构版本
		var archInUrl string
		switch arch {
		case "x86_64":
			archInUrl = "x64"
		case "aarch64":
			archInUrl = "arm64"
		default:
			return fmt.Errorf("不支持的系统架构: %s", arch)
		}

		g.SoftwareLogger.Infof("检测到系统架构: %s", arch)
		g.SoftwareLogger.Infof("开始下载 GitHub Runner...")
		downloadUrl := fmt.Sprintf("https://github.com/actions/runner/releases/download/v%s/actions-runner-linux-%s-%s.tar.gz", version, archInUrl, version)

		err = g.Download(downloadUrl, g.targetDir+"/actions-runner-linux.tar.gz")
		if err != nil {
			return fmt.Errorf("下载 runner 失败: %s", err)
		}

		// 解压
		g.SoftwareLogger.Infof("开始解压 GitHub Runner...")
		err, _ = g.RunShell("tar", "xzf", g.targetDir+"/actions-runner-linux.tar.gz", "-C", g.targetDir)
		if err != nil {
			return fmt.Errorf("解压失败: %s", err)
		}

		// 删除压缩包
		g.SoftwareLogger.Infof("删除压缩包，因为已解压到 " + g.targetDir)
		err = os.Remove(g.targetDir + "/actions-runner-linux.tar.gz")
		if err != nil {
			return fmt.Errorf("删除压缩包失败: %s", err)
		}

		// 安装依赖
		g.SoftwareLogger.Infof("开始安装依赖...")
		err, _ = g.RunShell("/bin/bash", g.targetDir+"/bin/installdependencies.sh")
		if err != nil {
			return fmt.Errorf("安装依赖失败: %s", err)
		}

		// 修改目录所有权
		g.SoftwareLogger.Infof("修改目录所有权...")
		err, _ = g.RunShell("chown", "-R", runnerUser+":"+runnerUser, g.targetDir)
		if err != nil {
			return fmt.Errorf("修改目录所有权失败: %s", err)
		}

		g.SoftwareLogger.Success("GitHub Runner 安装完成")
		return nil

	default:
		errMsg := fmt.Sprintf("不支持的操作系统: %s", osType)
		return g.SoftwareLogger.LogAndReturnErrorf("%s", errMsg)
	}
}
