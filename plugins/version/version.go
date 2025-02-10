package version

import (
	"servon/core"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Version 结构体通过嵌入core来继承其所有方法
type Version struct {
	*core.Core // 嵌入core，这样Version就继承了Core的所有方法
}

// 创建新的Version实例
func New(core *core.Core) *Version {
	return &Version{
		Core: core, // 直接将core赋值给嵌入字段
	}
}

func Setup(core *core.Core) {
	v := New(core)

	core.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "显示版本信息",
		Long:  "显示版本信息",
		Run: func(cmd *cobra.Command, args []string) {
			color.Green("Servon 版本: %s", v.GetVersion())
		},
	})
}
