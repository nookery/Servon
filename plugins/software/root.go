package software

import (
	"servon/core"

	"github.com/spf13/cobra"
)

type SoftWarePlugin struct {
	core *core.Core
}

func NewSoftWarePlugin(core *core.Core) *SoftWarePlugin {
	return &SoftWarePlugin{
		core: core,
	}
}

// Setup 注册到内核
func Setup(core *core.Core) error {
	plugin := NewSoftWarePlugin(core)
	core.AddCommand(plugin.GetSoftwareCommand())
	return nil
}

// GetSoftwareCommand 返回 software 命令
func (p *SoftWarePlugin) GetSoftwareCommand() *cobra.Command {
	cmd := p.core.NewCommand(core.CommandOptions{
		Use:   "software",
		Short: "软件管理",
	})

	cmd.AddCommand(p.newListCmd())
	cmd.AddCommand(p.newInstallCmd())
	cmd.AddCommand(p.newInfoCmd())
	cmd.AddCommand(p.newStartCmd())
	cmd.AddCommand(p.newStopCmd())

	return cmd
}
