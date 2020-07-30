package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"guguyun/models"
	"guguyun/util"
	"io"
	"log"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"
)

// 有关file结构体的业务逻辑

var sharePath map[int]string

// Upload 上传文件
func Upload(c *gin.Context) (res bool, msg string, err error) {
	file, err := c.FormFile("fileName")
	if err != nil {
		log.Printf("read file failed: %v", err)
		return false, "读取文件失败！", err
	}
	uploader, exists := c.Get("username")
	if !exists {
		log.Println("username doesn't exist!")
		return false, "登录认证已过期，请重新登录！", nil
	}
	destination := c.PostForm("path")
	dst := path.Join("./file/", uploader.(string), destination)
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		log.Printf("save file failed: %v", err)
		return false, "保存文件失败！", err
	}
	res = models.Upload(uploader.(string), file.Filename, dst)
	if !res {
		log.Println("file exists!")
		return false, "该路径下存在相同文件！", nil
	}
	return true, "上传文件成功！", nil
}

// Download 下载文件
func Download(c *gin.Context) (res bool, msg string, err error, file models.File) {
	filePath := c.Query("path")
	filePath = strings.ReplaceAll(filePath, "%2F", "/")
	filePath = path.Join("./file/", filePath)

	username, exists := c.Get("username")
	if !exists {
		log.Println("username doesn't exist!")
		return false, "登录认证已过期，请重新登录！", nil, models.File{}
	}

	res, file = models.GetFile(filePath)
	if !res {
		log.Println("file doesn't exist!")
		return false, "文件不存在！", err, models.File{}
	}

	if params := strings.Split(filePath, "/"); params[2] != username || !file.ShareAble {
		log.Println("not authorized!")
		return false, "没有权限！", nil, models.File{}
	}

	var open *os.File
	open, err = os.Open(filePath)
	if err != nil {
		log.Printf("open file failed:%v", err)
		return false, "文件不存在！", err, models.File{}
	}
	defer open.Close()

	fileName := path.Base(filePath)
	c.Writer.Header().Add("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	// c.File(filePath)
	_, err = io.Copy(c.Writer, open)
	if err != nil {
		log.Printf("io.copy failed:%v", err)
		return false, "下载文件失败！", nil, models.File{}
	}
	return true, "读取文件成功！", nil, file
}

// ChangeInformation 修改文件信息
func ChangeInformation(c *gin.Context) (res bool, msg string) {
	filePath := c.Query("path")
	result, file := models.GetFile(filePath)
	if !result {
		return false, "文件不存在！"
	}
	newName := c.PostForm("newName")
	newPath := c.PostForm("newPath")
	if newName == "" {
		newName = file.FileName
	}
	if newPath == "" {
		newPath = filePath
	}
	models.ChangeInformation(file, newPath, newName)
	return true, "修改文件信息成功！"
}

// DeleteFile 删除文件
func DeleteFile(c *gin.Context) bool {
	filePath := c.Query("path")
	err := os.Remove(filePath)
	if err != nil {
		return false
	}
	return models.DeleteFile(filePath)
}

// GetFiles 根据文件名或路径查找文件信息切片
func GetFiles(c *gin.Context) (bool, string, error, interface{}) {
	uploader, exists := c.Get("username")
	if !exists {
		log.Println("username does't exist!")
		return false, "签证已过期，请重新登陆！", nil, nil
	}
	fileName := c.Query("fileName")
	filePath := c.Query("path")
	files := models.GetFiles(uploader.(string), fileName, filePath)
	bytes, err := json.Marshal(&files)
	if err != nil {
		log.Printf("marshal failed:%v", err)
		return false, "序列化失败！", err, nil
	}
	return true, "搜索结果如下", nil, string(bytes)
}

// GetAllFiles 根据鉴权签证信息查找所有文件信息切片
func GetAllFiles(c *gin.Context) (bool, string, error, interface{}) {
	uploader, exists := c.Get("username")
	if !exists {
		log.Println("username doesn't exist!")
		return false, "签证已过期，请重新登录", nil, nil
	}
	files := models.GetAllFiles(uploader.(string))
	bytes, err := json.Marshal(&files)
	if err != nil {
		log.Printf("marshal failed:%v", err)
		return false, "序列化失败！", err, nil
	}
	return true, "所有文件如下", nil, string(bytes)
}


// Share 分享文件
func Share(s *gin.Context) (bool, string, error, string) {
	filePath := s.Query("path")
	res, file := models.GetFile(filePath)
	if !res {
		log.Println("file doesn't exist")
		return false, "文件不存在！", nil, ""
	}
	models.Share(file)
	// 利用随机数生成key键值，即邀请码
	rand.Seed(time.Now().UnixNano())
	key := rand.Intn(8999)
	for ; sharePath[key] == ""; {
		key = rand.Intn(8999)
	}
	sharePath[key] = filePath
	// 通过Encrypt函数加密
	rePath := util.Encrypt(filePath)
	scheme := s.Request.URL.Scheme
	host := s.Request.URL.Host
	url := path.Join(scheme, "://", host, ":8080", "/s/?value=", rePath)
	// 通过WriteFile函数生成二维码
	err := qrcode.WriteFile(url, qrcode.Medium, 256, "./qrcode/"+file.FileName+".png")
	if err != nil {
		log.Printf("get qrCode failed:%v", err)
		return false, "二维码生成失败！", err, ""
	}
	// 获取二维码图片
	open, err := os.Open("./qrcode/" + file.FileName + ".png")
	if err != nil {
		log.Printf("open file failed:%v", err)
		return false, "打开二维码文件失败！", err, ""
	}
	defer open.Close()
	// 传出信息
	result := map[string]interface{}{"key": key, "url": url, "qrCodePath": "./qrcode/" + file.FileName + ".png"}
	bytes, err := json.Marshal(&result)
	if err != nil {
		log.Printf("marshal failed:%v", err)
		return false, "序列化失败！", err, ""
	}
	return true, "分享链接已生成！", nil, string(bytes)
}

/*
// BreakPointRenewal 断点续传
func BreakPointRenewal(path string) {
	os.Open()
}
*/
