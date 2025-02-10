package contract

// SuperDeployMethod 定义部署操作的门面接口
type SuperDeployMethod interface {
	Deploy(logChan chan<- string) error
}
