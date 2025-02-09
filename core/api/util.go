package api

import (
	"servon/core/provider"
	"servon/core/utils"

	"github.com/spf13/cobra"
)

type Util struct {
	utilProvider provider.UtilProvider
}

func NewUtil() Util {
	return Util{
		utilProvider: provider.NewUtilProvider(),
	}
}
func (c *Util) PrintAndReturnError(errMsg string) error {
	return c.utilProvider.PrintAndReturnError(errMsg)
}

func (c *Util) PrintCommandHelp(cmd *cobra.Command) {
	utils.PrintCommandHelp(cmd)
}

func (c *Util) PrintCommandErrorAndExit(err error) error {
	return c.utilProvider.PrintCommandErrorAndExit(err)
}

func (c *Util) PrintCommandSuccess(msg string) {
	c.utilProvider.PrintCommandSuccess(msg)
}

func (c *Util) PrintStep(msg string) {
	c.utilProvider.PrintStep(msg)
}

func (c *Util) PrintStepSuccess(msg string) {
	c.utilProvider.PrintStepSuccess(msg)
}

func (c *Util) PrintSuccess(msg string) {
	c.utilProvider.PrintSuccess(msg)
}

func (c *Util) PrintStepFinish(msg string) {
	c.utilProvider.PrintStepFinish(msg)
}
