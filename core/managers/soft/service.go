package soft

import (
	"fmt"
	"servon/core/contract"
)

// ServiceManager 服务管理相关功能
type ServiceManager struct {
	*SoftManager
}

// RegisterService 注册后台服务软件
func (s *ServiceManager) RegisterService(name string, service contract.SuperService) error {
	s.LogUtil.Info("注册后台服务软件: " + name)

	if _, exists := s.Softwares[name]; exists {
		return fmt.Errorf("软件 %s 已注册为普通软件", name)
	}
	if _, exists := s.Services[name]; exists {
		return fmt.Errorf("软件 %s 已注册为后台服务软件", name)
	}

	s.Services[name] = service
	s.Softwares[name] = service
	return nil
}

// GetService 获取后台服务软件
func (s *ServiceManager) GetService(name string) (contract.SuperService, error) {
	service, ok := s.Services[name]
	if !ok {
		return nil, fmt.Errorf("后台服务软件 %s 未注册", name)
	}
	return service, nil
}

// GetAllServices 获取所有后台服务软件
func (s *ServiceManager) GetAllServices() []string {
	s.LogUtil.Info("获取所有后台服务软件...")
	serviceNames := make([]string, 0, len(s.Services))
	for name := range s.Services {
		serviceNames = append(serviceNames, name)
	}
	return serviceNames
}

// IsService 判断软件是否为后台服务软件
func (s *ServiceManager) IsService(name string) bool {
	_, ok := s.Services[name]
	return ok
}

// StartService 启动指定的后台服务
func (s *ServiceManager) StartService(name string) error {
	service, err := s.GetService(name)
	if err != nil {
		return err
	}
	return service.Start()
}

// StopService 停止指定的后台服务
func (s *ServiceManager) StopService(name string) error {
	service, err := s.GetService(name)
	if err != nil {
		return err
	}
	return service.Stop()
}

// RestartService 重启指定的后台服务
func (s *ServiceManager) RestartService(name string) error {
	if err := s.StopService(name); err != nil {
		return fmt.Errorf("停止服务失败: %v", err)
	}

	if err := s.StartService(name); err != nil {
		return fmt.Errorf("启动服务失败: %v", err)
	}

	return nil
}

// GetServiceStatus 获取指定后台服务的状态
func (s *ServiceManager) GetServiceStatus(name string) (map[string]string, error) {
	service, err := s.GetService(name)
	if err != nil {
		return nil, err
	}
	return service.GetStatus()
}
