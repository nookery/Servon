package core

import (
	"servon/core/utils"

	"github.com/spf13/cobra"
)

func (c *Core) PrintAndReturnError(errMsg string) error {
	return c.utilProvider.PrintAndReturnError(errMsg)
}

func (c *Core) PrintCommandHelp(cmd *cobra.Command) {
	utils.PrintCommandHelp(cmd)
}

func (c *Core) PrintCommandErrorAndExit(err error) error {
	return c.utilProvider.PrintCommandErrorAndExit(err)
}

func (c *Core) PrintCommandSuccess(msg string) {
	c.utilProvider.PrintCommandSuccess(msg)
}

func (c *Core) PrintStep(msg string) {
	c.utilProvider.PrintStep(msg)
}

func (c *Core) PrintStepSuccess(msg string) {
	c.utilProvider.PrintStepSuccess(msg)
}

func (c *Core) PrintSuccess(msg string) {
	c.utilProvider.PrintSuccess(msg)
}

func (c *Core) PrintStepFinish(msg string) {
	c.utilProvider.PrintStepFinish(msg)
}
