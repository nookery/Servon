package server

import (
	"os"

	"github.com/spf13/cobra"
)

// MakeRestartCommand 创建重启命令
func MakeRestartCommand() *cobra.Command {
	var verbose bool
	cmd := &cobra.Command{
		Use:   "restart",
		Short: "重启服务器",
		Run: func(cmd *cobra.Command, args []string) {
			verbose, _ = cmd.Flags().GetBool("verbose")
			server.SetVerbose(verbose)

			if err := server.RestartBackground(); err != nil {
				log.Error(err)
				os.Exit(1)
			}

			if !verbose {
				log.Success("服务器已重启 -> http://" + server.GetHost() + ":" + server.GetPortString())
			}
		},
	}

	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "启用详细日志模式，显示重启过程的详细信息")

	return cmd
}
