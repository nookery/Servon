package caddy

import (
	"servon/core"

	"github.com/spf13/cobra"
)

func (c *Caddy) NewInstallCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "安装 Caddy",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Install()
		},
	}
}

func (c *Caddy) NewProxyCommand() *cobra.Command {
	cmd := c.NewCommand(core.CommandOptions{
		Use:   "proxy",
		Short: "代理命令",
		Run: func(cmd *cobra.Command, args []string) {
			domain, _ := cmd.Flags().GetString("domain")
			target, _ := cmd.Flags().GetString("target")

			c.AddProxy(domain, target)
		},
	})

	cmd.Flags().String("domain", "", "域名")
	cmd.Flags().String("target", "", "目标地址")

	// 使用 Cobra 的自动验证
	cmd.MarkFlagRequired("domain")
	cmd.MarkFlagRequired("target")

	return cmd
}
