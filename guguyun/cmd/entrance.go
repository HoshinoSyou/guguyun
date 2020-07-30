package cmd

import (
	"github.com/gin-gonic/gin"
	"guguyun/controller"
	"guguyun/util/jwt"
)

// Entrance 存放路由
func Entrance() {
	r := gin.Default()
	acc := r.Group("/user")
	{
		acc.GET("/login", controller.Login)                   // 登录
		acc.POST("/register", controller.Register)            // 注册
		acc.Use(jwt.CheckToken())                             // jwt鉴权
		acc.PUT("/changePassword", controller.ChangePassword) // 修改密码
		acc.DELETE("/close", controller.CloseAccount)         // 注销账户/删除账户
	}
	file := r.Group("/file")
	{
		file.Use(jwt.CheckToken())
		file.GET("/download", controller.Download)                // 下载文件
		file.GET("/home", controller.GetAllFiles)                 // 显示自己的所有文件
		file.GET("/search", controller.GetFiles)                  // 搜索文件
		file.POST("/upload", controller.Upload)                   // 上传文件
		file.PUT("/changeName", controller.ChangeFileInformation) // 修改文件信息
		file.PUT("/share", controller.Share)                      // 分享文件
		file.DELETE("/delete", controller.DeleteFile)             // 删除文件
	}
	r.Run(":8080")
}
