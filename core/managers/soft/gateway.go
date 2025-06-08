package soft

import (
	"fmt"
	"servon/core/contract"
)

// GatewayManager 网关软件管理相关功能
type GatewayManager struct {
	*Manager
}

// RegisterGateway 注册网关软件
func (g *GatewayManager) RegisterGateway(name string, gateway contract.SuperGateway) error {
	g.LogUtil.Info("注册网关软件: " + name)

	if _, exists := g.Softwares[name]; exists {
		return fmt.Errorf("软件 %s 已注册为普通软件", name)
	}
	if _, exists := g.Gateways[name]; exists {
		return fmt.Errorf("软件 %s 已注册为网关软件", name)
	}

	g.Gateways[name] = gateway
	g.Softwares[name] = gateway
	return nil
}

// GetGateway 获取网关软件
func (g *GatewayManager) GetGateway(name string) (contract.SuperGateway, error) {
	gateway, ok := g.Gateways[name]
	if !ok {
		return nil, fmt.Errorf("网关软件 %s 未注册", name)
	}
	return gateway, nil
}

// GetAllGateways 获取所有网关软件
func (g *GatewayManager) GetAllGateways() []string {
	g.LogUtil.Info("获取所有网关软件...")
	gatewayNames := make([]string, 0, len(g.Gateways))
	for name := range g.Gateways {
		gatewayNames = append(gatewayNames, name)
	}
	return gatewayNames
}

// IsGateway 判断软件是否为网关软件
func (g *GatewayManager) IsGateway(name string) bool {
	_, ok := g.Gateways[name]
	return ok
}
