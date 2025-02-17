package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleListUsers 处理获取用户列表的请求
func (h *WebHandler) HandleListUsers(c *gin.Context) {
	users, err := h.App.UserLib.GetUserList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// HandleDeleteUser 处理删除用户的请求
func (h *WebHandler) HandleDeleteUser(c *gin.Context) {
	username := c.Param("username")
	err := h.App.UserLib.DeleteUser(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// HandleCreateUser 处理创建用户的请求
func (h *WebHandler) HandleCreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	err := h.App.UserLib.CreateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
