package server

import (
	"os"

	"github.com/spf13/cobra"
)

// MakeStopCommand 创建停止命令
func MakeStopCommand() *cobra.Command {
	var verbose bool

	cmd := &cobra.Command{
		Use:   "stop",
		Short: "停止服务器",
		Run: func(cmd *cobra.Command, args []string) {
			verbose, _ = cmd.Flags().GetBool("verbose")
			server.SetVerbose(verbose)

			if err := server.StopBackground(); err != nil {
				log.Error(err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "启用详细日志模式，显示停止过程的详细信息")

	return cmd
}
