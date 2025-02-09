package provider

import (
	"github.com/spf13/cobra"
)

// DataProvider 数据提供者
type DataProvider struct {
	RootCmd *cobra.Command
}

func NewDataProvider() DataProvider {
	return DataProvider{
		RootCmd: &cobra.Command{},
	}
}

// GetDataRootFolder 获取数据根文件夹
func (p *DataProvider) GetDataRootFolder() string {
	return "/data"
}
