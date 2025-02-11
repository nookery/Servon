package api

import (
	"servon/core/contract"
	"servon/core/libs"

	"github.com/spf13/cobra"
)

type DeployApi struct {
	methods map[string]contract.SuperDeployMethod
}

func NewDeployApi() DeployApi {
	return DeployApi{
		methods: make(map[string]contract.SuperDeployMethod),
	}
}

func (c *DeployApi) RegisterMethod(name string, method contract.SuperDeployMethod) {
	c.methods[name] = method
}

func (c *DeployApi) GetMethods() []string {
	methods := make([]string, 0, len(c.methods))
	for name := range c.methods {
		methods = append(methods, name)
	}
	return methods
}

func (c *DeployApi) GetDeployCommand() *cobra.Command {
	rootCmd := libs.NewCommand(libs.CommandOptions{
		Use:   "deploy",
		Short: "部署项目",
	})

	rootCmd.AddCommand(c.NewMethodsCommand())

	return rootCmd
}

func (c *DeployApi) NewMethodsCommand() *cobra.Command {
	return libs.NewCommand(libs.CommandOptions{
		Use:   "methods",
		Short: "列举出当前的部署方法",
		Run: func(cmd *cobra.Command, args []string) {
			printer.PrintList(c.GetMethods(), "部署方法")
		},
	})
}
