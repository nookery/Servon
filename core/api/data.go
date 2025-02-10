package api

import (
	"servon/core/libs"
)

type Data struct {
	dataProvider libs.DataProvider
}

func NewData() Data {
	return Data{
		dataProvider: libs.NewDataProvider(),
	}
}

func (c *Data) GetDataRootFolder() string {
	return c.dataProvider.GetDataRootFolder()
}
