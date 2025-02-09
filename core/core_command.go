package core

import "github.com/spf13/cobra"

// AddCommand 添加命令
func (c *Core) AddCommand(cmd *cobra.Command) {
	c.commandProvider.AddCommand(cmd)
}

// GetRootCommand 获取根命令
func (c *Core) GetRootCommand() *cobra.Command {
	return c.commandProvider.RootCmd
}
