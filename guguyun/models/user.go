package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"guguyun/dao"
)

// 有关user结构体的数据库（user表）增删改查相关操作实现

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register 注册 增加user记录
func Register(username string, password string) bool {
	dao.DB.AutoMigrate(&User{})
	var us User
	var u = User{
		Model:    gorm.Model{},
		Username: username,
		Password: password,
	}
	dao.DB.Where("Username = ?", username).First(&us)
	if us.ID > 0 {
		return false
	}
	dao.DB.Create(&u)
	return true
}

// Login 登录 查询user记录
func Login(username string, password string) bool {
	var u User
	fmt.Println(username,password)
	dao.DB.Where("username = ?", username).First(&u)
	if u.Password == password {
		return true
	} else {
		return false
	}
}

// ChangePassword 修改密码 修改user记录
func ChangePassword(username string, oldPassword string, newPassword string) (bool, string) {
	var u = User{
		Username: username,
		Password: oldPassword,
	}
	res := dao.DB.NewRecord(u)
	if !res {
		return false, "用户名或旧密码错误！"
	}
	dao.DB.Model(&u).Where("Username = ?").Update("Password", newPassword)
	return true, "修改密码成功！"
}

// CloseAccount 注销账户/删除账户 删除user记录
func CloseAccount(username string, password string) bool {
	var u User
	dao.DB.Where("Username = ?", username).First(&u)
	if u.Password == password {
		dao.DB.Delete(&u)
		return true
	} else {
		return false
	}
}
