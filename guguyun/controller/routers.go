package controller

import (
	"github.com/gin-gonic/gin"
	"guguyun/service"
	"guguyun/util/jwt"
	"guguyun/util/response"
)

// 控制器
// 路由方法的实现

// Login 登录
func Login(l *gin.Context) {
	res := service.Login(l)
	if res {
		username := l.PostForm("username")
		token := jwt.Create(username)
		response.OkWithToken(l, "登陆成功！欢迎回来！", token)
	} else {
		response.Error(l, "登录失败，用户名或密码错误！", nil)
	}
}

// Register 注册
func Register(r *gin.Context) {
	res := service.Register(r)
	if res {
		response.Ok(r, "注册成功！")
	} else {
		response.Error(r, "注册失败，用户名已存在！", nil)
	}
}

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
	res, msg := service.ChangePassword(c)
	if res {
		response.Ok(c, msg)
	} else {
		response.Error(c, msg, nil)
	}
}

// CloseAccount 注销账户/删除账户
func CloseAccount(a *gin.Context) {
	res := service.CloseAccount(a)
	if res {
		response.Ok(a, "注销账户成功！")
	} else {
		response.Error(a, "注销账户失败，用户名或密码错误！", nil)
	}
}

// Upload 上传文件
func Upload(u *gin.Context) {
	res, msg, err := service.Upload(u)
	if res {
		response.Ok(u, msg)
	} else {
		response.Error(u, msg, err)
	}
}

func Download(d *gin.Context) {
	res, msg, err, file := service.Download(d)
	if res {
		response.OkWithData(d, msg, file)
	} else {
		response.Error(d, msg, err)
	}
}

func ChangeFileInformation(c *gin.Context) {
	res, msg := service.ChangeInformation(c)
	if res {
		response.Ok(c, msg)
	} else {
		response.Error(c, msg, nil)
	}
}

func DeleteFile(d *gin.Context) {
	res := service.DeleteFile(d)
	if res {
		response.Ok(d, "删除文件成功！")
	} else {
		response.Error(d, "删除文件失败，文件不存在！", nil)
	}
}

func GetFiles(g *gin.Context) {
	res, msg, err, files := service.GetFiles(g)
	if res {
		response.OkWithData(g, msg, files)
	} else {
		response.Error(g, msg, err)
	}
}

func GetAllFiles(g *gin.Context) {
	res, msg, err, data := service.GetAllFiles(g)
	if res {
		response.OkWithData(g, msg, data)
	} else {
		response.Error(g, msg, err)
	}
}

func Share(s *gin.Context) {
	res, msg, err, data := service.Share(s)
	if res {
		response.OkWithData(s, msg, data)
	} else {
		response.Error(s, msg, err)
	}
}
