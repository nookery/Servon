package core

func (c *Core) RunShell(command string, args ...string) error {
	return c.shellProvider.Execute(command, args...)
}

func (c *Core) RunShellWithOutput(command string, args ...string) (string, error) {
	return c.shellProvider.ExecuteWithOutput(command, args...)
}
