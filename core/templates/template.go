package templates

import (
	"embed"
	"fmt"

	"github.com/fatih/color"
)

//go:embed systemd_service.tmpl
var SystemdServiceTemplateFS embed.FS

//go:embed supervisor.conf.tmpl
var SupervisorConfigTemplateFS embed.FS

var (
	cyan   = color.New(color.FgCyan).SprintFunc()
	purple = color.New(color.FgMagenta).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	blue   = color.New(color.FgBlue).SprintFunc()
)

// UsageTemplate 返回命令使用说明模板
func UsageTemplate() string {
	return `
` + purple("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━") + `
` + yellow("📌 命令:") + ` {{.UseLine}}
` + green("📝 描述:") + ` {{.Short}}

` + blue("🎯 参数列表:") + `
{{.LocalFlags.FlagUsages}}
` + cyan("✨ 示例:") + `{{.CommandPath}} [参数]

` + yellow("💡 提示:") + ` 使用 -h 或 --help 查看更多帮助信息
` + purple("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━") + `
`
}

// GetSupervisorConfigTemplate 返回supervisor配置模板内容
func GetSupervisorConfigTemplate() (string, error) {
	tmplContent, err := SupervisorConfigTemplateFS.ReadFile("supervisor.conf.tmpl")
	if err != nil {
		return "", fmt.Errorf("读取supervisor模板文件失败: %v", err)
	}
	return string(tmplContent), nil
}
