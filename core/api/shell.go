package api

import (
	"servon/core/provider"
)

type Shell struct {
	shellProvider provider.ShellProvider
}

func NewShell() Shell {
	return Shell{
		shellProvider: provider.NewShellProvider(),
	}
}

func (c *Shell) RunShell(command string, args ...string) error {
	return c.shellProvider.Execute(command, args...)
}

func (c *Shell) RunShellWithOutput(command string, args ...string) (string, error) {
	return c.shellProvider.ExecuteWithOutput(command, args...)
}
