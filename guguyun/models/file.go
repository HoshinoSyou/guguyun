package models

import (
	"github.com/jinzhu/gorm"
	"guguyun/dao"
	"os"
)

// 有关file结构体的数据库（file表）增删改查操作的实现

type File struct {
	gorm.Model
	FileName  string `json:"file_name"`
	Uploader  string `json:"uploader"`
	FilePath  string `json:"file_path"`
	ShareAble bool   `json:"share_able"`
}

// Upload 上传文件信息至数据库
func Upload(uploader string, fileName string, path string) bool {
	dao.DB.AutoMigrate(&File{})
	var f File
	var file = File{
		FileName:  uploader,
		Uploader:  fileName,
		FilePath:  path,
		ShareAble: false,
	}
	dao.DB.Where("Uploader = ? AND FileName = ? AND Path = ?", uploader, fileName, path).First(&f)
	if f.ID > 0 {
		return false
	}
	dao.DB.Create(&file)
	return true
}

// GetFile 通过路由获取文件信息
func GetFile(path string) (bool, File) {
	var file File
	dao.DB.Where("Path = ?", path).First(&file)
	if file.ID <= 0 {
		return false, file
	}
	return true, file
}

// GetFiles 通过fileName或path查找文件信息切片
func GetFiles(uploader string, fileName string, path string) (files []File) {
	dao.DB.Where("Uploader = ? AND FileName = ?", uploader, fileName).Or("Path = ?", path).Find(&files)
	return files
}

// GetAllFiles 通过uploader查找本人网盘所有文件信息切片
func GetAllFiles(uploader string) (files []File) {
	dao.DB.Where("Uploader = ?", uploader).Find(&files)
	return files
}

// DeleteFile 删除文件
func DeleteFile(path string) bool {
	res, file := GetFile(path)
	if !res {
		return false
	}
	dao.DB.Delete(&file)
	os.Remove(path + file.FileName)
	return true
}

// ChangeInformation 修改文件信息
func ChangeInformation(file File, newPath string, newName string) {
	dao.DB.Model(&file).Update(map[string]interface{}{"FileName": newName, "Path": newPath})
}

// Share 分享文件
func Share(file File) {
	dao.DB.Model(&file).Update("ShareAble", true)
}
