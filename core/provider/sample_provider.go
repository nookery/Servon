package provider

import (
	"github.com/spf13/cobra"
)

// SampleProvider 示例提供者
type SampleProvider struct {
	RootCmd *cobra.Command
}

func NewSampleProvider() SampleProvider {
	return SampleProvider{
		RootCmd: &cobra.Command{},
	}
}
