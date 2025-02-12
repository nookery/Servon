package caddy

import (
	"servon/core"

	"github.com/spf13/cobra"
)

func NewProxyCommand(caddy *Caddy) *cobra.Command {
	cmd := caddy.NewCommand(core.CommandOptions{
		Use:   "proxy",
		Short: "代理命令",
		Run: func(cmd *cobra.Command, args []string) {
			domain, _ := cmd.Flags().GetString("domain")
			target, _ := cmd.Flags().GetString("target")

			caddy.AddProxy(domain, target)
		},
	})

	cmd.Flags().String("domain", "", "域名")
	cmd.Flags().String("target", "", "目标地址")

	// 使用 Cobra 的自动验证
	cmd.MarkFlagRequired("domain")
	cmd.MarkFlagRequired("target")

	return cmd
}
