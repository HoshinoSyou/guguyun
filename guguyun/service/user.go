package service

import (
	"github.com/gin-gonic/gin"
	"guguyun/models"
)

/*
有关user结构体的业务逻辑
 */

// Login 登录
func Login(c *gin.Context) bool {
	username := c.PostForm("username")
	password := c.PostForm("password")
	return models.Login(username, password)
}

// Register 注册
func Register(c *gin.Context) bool {
	username := c.PostForm("username")
	password := c.PostForm("password")
	return models.Register(username, password)
}

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) (bool, string) {
	username := c.PostForm("username")
	oldPassword := c.PostForm("oldPassword")
	newPassword := c.PostForm("newPassword")
	return models.ChangePassword(username, oldPassword, newPassword)
}

// CloseAccount 注销账户/删除账户
func CloseAccount(c *gin.Context) bool {
	username := c.PostForm("username")
	password := c.PostForm("password")
	return models.CloseAccount(username, password)
}
