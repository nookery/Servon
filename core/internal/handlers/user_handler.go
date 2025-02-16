package handlers

import (
	"net/http"
	"servon/core/internal/libs"

	"github.com/gin-gonic/gin"
)

var UserManager = libs.DefaultUserManager

// HandleListUsers 处理获取用户列表的请求
func HandleListUsers(c *gin.Context) {
	users, err := UserManager.GetUserList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// HandleDeleteUser 处理删除用户的请求
func HandleDeleteUser(c *gin.Context) {
	username := c.Param("username")
	err := UserManager.DeleteUser(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// HandleCreateUser 处理创建用户的请求
func HandleCreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	err := UserManager.CreateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
