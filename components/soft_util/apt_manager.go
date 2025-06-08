package soft_util

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// AptPackage 表示一个 apt 包的信息
type AptPackage struct {
	Name          string    `json:"name"`
	Version       string    `json:"version"`
	Architecture  string    `json:"architecture"`
	Status        string    `json:"status"`
	InstalledSize int64     `json:"installed_size"`
	Description   string    `json:"description"`
	InstallDate   time.Time `json:"install_date,omitempty"`
	IsInstalled   bool      `json:"is_installed"`
}

// ShellUtil 定义了执行shell命令的接口
type ShellUtil interface {
	RunShell(command string, args ...string) (error, string)
	RunShellWithSudo(command string, args ...string) (error, string)
}

// LogUtil 定义了日志记录的接口
type LogUtil interface {
	Info(message string)
	Warn(message string)
	LogAndReturnErrorf(format string, args ...interface{}) error
}

// AptManager 管理APT包操作
type AptManager struct {
	ShellUtil ShellUtil
	LogUtil   LogUtil
}

// NewAptManager 创建一个新的APT管理器实例
func NewAptManager(shellUtil ShellUtil, logUtil LogUtil) *AptManager {
	return &AptManager{
		ShellUtil: shellUtil,
		LogUtil:   logUtil,
	}
}

// AptUpdate 更新软件包索引
func (p *AptManager) AptUpdate() (string, error) {
	err, output := p.ShellUtil.RunShellWithSudo("apt-get", "update")
	if err != nil {
		return output, p.LogUtil.LogAndReturnErrorf("更新索引失败: %v", err)
	}

	p.LogUtil.Info("Apt 索引更新成功")
	return output, nil
}

// AptInstall 安装指定的软件包
func (p *AptManager) AptInstall(packages ...string) error {
	if len(packages) == 0 {
		return p.LogUtil.LogAndReturnErrorf("未指定要安装的软件包")
	}

	// 使用数组传递参数，避免命令注入
	args := append([]string{"install", "-y"}, packages...)
	err, output := p.ShellUtil.RunShellWithSudo("apt-get", args...)
	if err != nil {
		return p.LogUtil.LogAndReturnErrorf("安装失败: %v, 输出: %s", err, output)
	}

	p.LogUtil.Info(fmt.Sprintf("Apt 安装成功: %v", packages))
	return nil
}

// AptRemove 移除指定的软件包
func (p *AptManager) AptRemove(packages ...string) error {
	if len(packages) == 0 {
		return p.LogUtil.LogAndReturnErrorf("未指定要移除的软件包")
	}

	// 使用数组传递参数，避免命令注入
	args := append([]string{"remove", "-y"}, packages...)
	err, output := p.ShellUtil.RunShellWithSudo("apt-get", args...)
	if err != nil {
		return p.LogUtil.LogAndReturnErrorf("移除失败: %v, 输出: %s", err, output)
	}

	p.LogUtil.Info(fmt.Sprintf("Apt 移除成功: %v", packages))
	return nil
}

// AptPurge 完全移除软件包及其配置文件
func (p *AptManager) AptPurge(packages ...string) error {
	if len(packages) == 0 {
		return p.LogUtil.LogAndReturnErrorf("未指定要清理的软件包")
	}

	// 使用数组传递参数，避免命令注入
	args := append([]string{"purge", "-y"}, packages...)
	err, output := p.ShellUtil.RunShellWithSudo("apt-get", args...)
	if err != nil {
		return p.LogUtil.LogAndReturnErrorf("清理失败: %v, 输出: %s", err, output)
	}

	p.LogUtil.Info(fmt.Sprintf("Apt 清理成功: %v", packages))
	return nil
}

// AptIsInstalled 检查软件包是否已安装
func (p *AptManager) AptIsInstalled(packageName string) bool {
	err, output := p.ShellUtil.RunShell("dpkg-query", "-W", "-f=${Status}", packageName)
	if err != nil {
		return false
	}

	return strings.Contains(output, "install ok installed")
}

// AptGetPackageInfo 获取软件包详细信息
func (p *AptManager) AptGetPackageInfo(packageName string) (*AptPackage, error) {
	if !p.AptIsInstalled(packageName) {
		return nil, fmt.Errorf("软件包 %s 未安装", packageName)
	}

	// 获取包信息
	err, output := p.ShellUtil.RunShell("dpkg-query", "-W", "-f=${Package}|${Version}|${Architecture}|${Status}|${Installed-Size}|${Description}\\n", packageName)
	if err != nil {
		return nil, p.LogUtil.LogAndReturnErrorf("获取软件包信息失败: %v", err)
	}

	parts := strings.Split(strings.TrimSpace(output), "|")
	if len(parts) < 6 {
		return nil, fmt.Errorf("无法解析软件包信息: %s", output)
	}

	var installedSize int64
	if size, err := strconv.ParseInt(parts[4], 10, 64); err == nil {
		installedSize = size
	}

	// 获取安装时间
	var installDate time.Time
	err, dateOutput := p.ShellUtil.RunShell("stat", "-c", "%Y", fmt.Sprintf("/var/lib/dpkg/info/%s.list", packageName))
	if err == nil {
		if timestamp, err := strconv.ParseInt(strings.TrimSpace(dateOutput), 10, 64); err == nil {
			installDate = time.Unix(timestamp, 0)
		}
	}

	return &AptPackage{
		Name:          parts[0],
		Version:       parts[1],
		Architecture:  parts[2],
		Status:        parts[3],
		InstalledSize: installedSize,
		Description:   parts[5],
		InstallDate:   installDate,
		IsInstalled:   true,
	}, nil
}

// AptListInstalled 列出所有已安装的软件包
func (p *AptManager) AptListInstalled() ([]AptPackage, error) {
	err, output := p.ShellUtil.RunShell("dpkg-query", "-W", "-f=${Package}\\n")
	if err != nil {
		return nil, p.LogUtil.LogAndReturnErrorf("获取已安装软件包列表失败: %v", err)
	}

	packages := []AptPackage{}
	packageNames := strings.Split(strings.TrimSpace(output), "\n")

	for _, name := range packageNames {
		if name == "" {
			continue
		}

		info, err := p.AptGetPackageInfo(name)
		if err != nil {
			p.LogUtil.Warn(fmt.Sprintf("获取软件包 %s 信息失败: %v", name, err))
			continue
		}

		packages = append(packages, *info)
	}

	return packages, nil
}

// AptSearch 搜索软件包
func (p *AptManager) AptSearch(keyword string) ([]AptPackage, error) {
	err, output := p.ShellUtil.RunShell("apt-cache", "search", keyword)
	if err != nil {
		return nil, p.LogUtil.LogAndReturnErrorf("搜索软件包失败: %v", err)
	}

	results := []AptPackage{}
	lines := strings.Split(strings.TrimSpace(output), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, " - ", 2)
		if len(parts) < 2 {
			continue
		}

		packageName := strings.TrimSpace(parts[0])
		description := strings.TrimSpace(parts[1])

		// 获取版本信息
		err, versionOutput := p.ShellUtil.RunShell("apt-cache", "policy", packageName)
		version := ""
		if err == nil {
			// 使用正则表达式提取版本信息
			re := regexp.MustCompile(`Candidate:\s+(\S+)`)
			matches := re.FindStringSubmatch(versionOutput)
			if len(matches) > 1 {
				version = matches[1]
			}
		}

		results = append(results, AptPackage{
			Name:        packageName,
			Version:     version,
			Description: description,
			IsInstalled: p.AptIsInstalled(packageName),
		})
	}

	return results, nil
}

// AptUpgrade 升级所有软件包
func (p *AptManager) AptUpgrade(distUpgrade bool) (string, error) {
	var command string
	if distUpgrade {
		command = "dist-upgrade"
	} else {
		command = "upgrade"
	}

	err, output := p.ShellUtil.RunShellWithSudo("apt-get", command, "-y")
	if err != nil {
		return "", p.LogUtil.LogAndReturnErrorf("升级失败: %v", err)
	}

	p.LogUtil.Info("Apt 升级成功")
	return output, nil
}

// AptAutoRemove 自动移除不再需要的依赖包
func (p *AptManager) AptAutoRemove() error {
	err, _ := p.ShellUtil.RunShellWithSudo("apt-get", "autoremove", "-y")
	if err != nil {
		return p.LogUtil.LogAndReturnErrorf("自动移除失败: %v", err)
	}

	p.LogUtil.Info("Apt 自动移除成功")
	return nil
}

// AptClean 清理本地缓存的软件包
func (p *AptManager) AptClean() error {
	err, _ := p.ShellUtil.RunShellWithSudo("apt-get", "clean")
	if err != nil {
		return p.LogUtil.LogAndReturnErrorf("清理缓存失败: %v", err)
	}

	p.LogUtil.Info("Apt 缓存清理成功")
	return nil
}

// AptGetUpgradable 获取可升级的软件包列表
func (p *AptManager) AptGetUpgradable() ([]AptPackage, error) {
	err, output := p.ShellUtil.RunShell("apt", "list", "--upgradable")
	if err != nil {
		return nil, p.LogUtil.LogAndReturnErrorf("获取可升级软件包列表失败: %v", err)
	}

	upgradable := []AptPackage{}
	lines := strings.Split(output, "\n")

	// 跳过第一行（标题行）
	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		// 解析输出格式: name/source arch [version] [repo]
		parts := strings.Split(line, "/")
		if len(parts) < 2 {
			continue
		}

		name := parts[0]

		// 提取版本信息
		versionMatch := regexp.MustCompile(`\[(.*?)\]`).FindStringSubmatch(line)
		version := ""
		if len(versionMatch) > 1 {
			version = versionMatch[1]
		}

		// 提取架构信息
		archMatch := regexp.MustCompile(`\s+(\S+)\s+\[`).FindStringSubmatch(line)
		arch := ""
		if len(archMatch) > 1 {
			arch = archMatch[1]
		}

		upgradable = append(upgradable, AptPackage{
			Name:         name,
			Version:      version,
			Architecture: arch,
			IsInstalled:  true,
		})
	}

	return upgradable, nil
}

// AptGetDependencies 获取软件包的依赖关系
func (p *AptManager) AptGetDependencies(packageName string) ([]string, error) {
	err, output := p.ShellUtil.RunShell("apt-cache", "depends", packageName)
	if err != nil {
		return nil, p.LogUtil.LogAndReturnErrorf("获取依赖关系失败: %v", err)
	}

	dependencies := []string{}
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Depends:") || strings.HasPrefix(line, "PreDepends:") {
			dep := strings.TrimSpace(strings.Split(line, ":")[1])
			if dep != "" && !strings.Contains(dep, "<") && !strings.Contains(dep, ">") {
				dependencies = append(dependencies, dep)
			}
		}
	}

	return dependencies, nil
}

// AptGetReverseDependencies 获取依赖于指定软件包的包列表
func (p *AptManager) AptGetReverseDependencies(packageName string) ([]string, error) {
	err, output := p.ShellUtil.RunShell("apt-cache", "rdepends", packageName)
	if err != nil {
		return nil, p.LogUtil.LogAndReturnErrorf("获取反向依赖关系失败: %v", err)
	}

	dependencies := []string{}
	lines := strings.Split(output, "\n")

	// 跳过前两行（标题行）
	for i := 2; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line != "" && !strings.HasPrefix(line, "Reverse") {
			dependencies = append(dependencies, line)
		}
	}

	return dependencies, nil
}

// AptToJSON 将软件包信息转换为JSON字符串
func (p *AptManager) AptToJSON(packages []AptPackage) (string, error) {
	data, err := json.MarshalIndent(packages, "", "  ")
	if err != nil {
		return "", fmt.Errorf("转换为JSON失败: %v", err)
	}

	return string(data), nil
}
