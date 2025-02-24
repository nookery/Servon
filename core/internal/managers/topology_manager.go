package managers

import (
	"servon/core/internal/contract"
)

type TopologyManager struct {
	softManager *SoftManager
}

func NewTopologyManager(softManager *SoftManager) *TopologyManager {
	return &TopologyManager{
		softManager: softManager,
	}
}

func (m *TopologyManager) GetAllGateways() []string {
	return m.softManager.GetAllGateways()
}

func (m *TopologyManager) GetProjects(gatewayName string) ([]contract.Project, error) {
	gateway, err := m.softManager.GetGateway(gatewayName)
	if err != nil {
		return nil, err
	}
	return gateway.GetProjects()
}

func (m *TopologyManager) AddProject(gatewayName string, project contract.Project) error {
	gateway, err := m.softManager.GetGateway(gatewayName)
	if err != nil {
		return err
	}
	return gateway.AddProject(project)
}

func (m *TopologyManager) RemoveProject(gatewayName string, projectName string) error {
	gateway, err := m.softManager.GetGateway(gatewayName)
	if err != nil {
		return err
	}
	return gateway.RemoveProject(projectName)
}

// GetGateway 获取指定的网关软件
func (m *TopologyManager) GetGateway(gatewayName string) (contract.SuperGateway, error) {
	gateway, err := m.softManager.GetGateway(gatewayName)
	if err != nil {
		return nil, err
	}
	return gateway, nil
}
