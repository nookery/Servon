package templates

import "github.com/fatih/color"

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
