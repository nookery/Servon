package api

import (
	"servon/core/provider"
)

type Data struct {
	dataProvider provider.DataProvider
}

func NewData() Data {
	return Data{
		dataProvider: provider.NewDataProvider(),
	}
}

func (c *Data) GetDataRootFolder() string {
	return c.dataProvider.GetDataRootFolder()
}
