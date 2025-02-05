package cmd

import (
	"servon/cmd/software"

	"github.com/spf13/cobra"
)

// SoftwareCmd 表示 software 命令
var SoftwareCmd = &cobra.Command{}

func init() {
	SoftwareCmd = software.GetSoftwareCommand()
}
