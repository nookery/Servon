package api

import (
	"servon/core/provider"

	"github.com/spf13/cobra"
)

type Command struct {
	commandProvider provider.CommandProvider
}

func NewCommand() Command {
	return Command{
		commandProvider: provider.NewCommandProvider(),
	}
}

func (c *Command) AddCommand(cmd *cobra.Command) {
	c.commandProvider.AddCommand(cmd)
}

func (c *Command) GetRootCommand() *cobra.Command {
	return c.commandProvider.RootCmd
}
