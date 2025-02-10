package caddy

import (
	"github.com/spf13/cobra"
)

func NewProxyCommand(caddy *Caddy) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proxy",
		Short: "代理命令",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Flags().Parse(args); err != nil {
				cmd.PrintErrf("参数错误: %v\n", err)
				cmd.PrintErrf("用法示例: servon caddy proxy --domain example.com --target http://localhost:8080\n")
				return err
			}
			domain, _ := cmd.Flags().GetString("domain")
			target, _ := cmd.Flags().GetString("target")
			caddy.AddProxy(domain, target)
			return nil
		},
	}

	cmd.Flags().String("domain", "", "域名")
	cmd.Flags().String("target", "", "目标地址")

	// 标记必需的参数
	cmd.MarkFlagRequired("domain")
	cmd.MarkFlagRequired("target")

	return cmd
}
