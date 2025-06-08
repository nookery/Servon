package xcode_util

// BuildTargetInfo 构建目标信息结构体
type BuildTargetInfo struct {
	ProjectFile     string // 项目文件路径
	ProjectType     string // 项目类型 ("workspace" 或 "project")
	ProjectTypeName string // 项目类型显示名称
	Scheme          string // 构建方案
	ProjectArchs    string // 项目支持的架构
	TargetArch      string // 构建目标架构
}
