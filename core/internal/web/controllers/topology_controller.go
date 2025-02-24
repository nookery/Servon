package controllers

import (
	"servon/core/internal/contract"
	"servon/core/internal/managers"

	"github.com/gin-gonic/gin"
)

type TopologyController struct {
	topologyManager *managers.TopologyManager
}

func NewTopologyController(manager *managers.TopologyManager) *TopologyController {
	return &TopologyController{
		topologyManager: manager,
	}
}

func (c *TopologyController) HandleGetGateways(ctx *gin.Context) {
	gateways := c.topologyManager.GetAllGateways()
	ctx.JSON(200, gateways)
}

func (c *TopologyController) HandleGetProjects(ctx *gin.Context) {
	gateway := ctx.Param("gateway")
	projects, err := c.topologyManager.GetProjects(gateway)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, projects)
}

func (c *TopologyController) HandleAddProject(ctx *gin.Context) {
	gateway := ctx.Param("gateway")
	var project contract.Project
	if err := ctx.ShouldBindJSON(&project); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := c.topologyManager.AddProject(gateway, project)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(200)
}

func (c *TopologyController) HandleRemoveProject(ctx *gin.Context) {
	gateway := ctx.Param("gateway")
	name := ctx.Param("name")

	err := c.topologyManager.RemoveProject(gateway, name)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(200)
}

func (c *TopologyController) HandleGetGatewayConfig(ctx *gin.Context) {
	gateway := ctx.Param("gateway")
	gatewayInstance, err := c.topologyManager.GetGateway(gateway)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	config, err := gatewayInstance.GetConfig()
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, config)
}

func (c *TopologyController) HandleSetGatewayConfig(ctx *gin.Context) {
	gateway := ctx.Param("gateway")
	var body struct {
		Config string `json:"config"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	gatewayInstance, err := c.topologyManager.GetGateway(gateway)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = gatewayInstance.SetConfig(map[string]interface{}{
		"config": body.Config,
	})
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(200)
}
