package core

import (
	"servon/core/libs"
)

// 调用关系 Core -> Core API(Public) -> Libs(Private)

import (
	"servon/core/api"
)

const DataRootFolder = "/data"
const LoggerFolder = "/logs"

type Core struct {
	api.Command
	api.Soft
	api.SystemApi
	api.Version
	api.Data
	api.Sample
	api.LogApi
}

type OSType = libs.OSType
type CommandOptions = api.CommandOptions

const (
	Ubuntu  OSType = "ubuntu"
	Debian  OSType = "debian"
	CentOS  OSType = "centos"
	RedHat  OSType = "redhat"
	Unknown OSType = "unknown"
)

// New 创建Core实例
func New() *Core {
	return &Core{
		Command:   api.NewCommandApi(),
		Soft:      api.NewSoft(),
		SystemApi: api.NewSystemApi(),
		Version:   api.NewVersion(),
		Data:      api.NewData(),
		Sample:    api.NewSample(),
		LogApi:    api.NewLogApi(),
	}
}
