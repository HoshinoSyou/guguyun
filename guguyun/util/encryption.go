package util

import (
	"encoding/base64"
	"log"
	"path"
	"strings"
)

// Encrypt 基于base64的加密函数
func Encrypt(path string) string {
	str := strings.Split(path, "/")
	i := len(str)
	var s string
	for j := 3; j < i-1; j++ {
		s = str[j] + "-" + s
	}
	word1 := base64.StdEncoding.EncodeToString([]byte(str[2] + "." + str[i-1]))
	word2 := base64.StdEncoding.EncodeToString([]byte(s))
	return word1 + "=" + word2
}

// UnEncrypt 基于base64的解密函数
func UnEncrypt(getPath string) (bool, string, error, error) {
	str := strings.Split(getPath, "=")
	str1, err1 := base64.StdEncoding.DecodeString(str[0])
	str2, err2 := base64.StdEncoding.DecodeString(str[1])
	if err1 != nil || err2 != nil {
		log.Printf("decode failed:%v,%v", err1, err2)
		return false, "", err1, err2
	}
	word1 := strings.Split(string(str1), ".")
	word2 := strings.Split(string(str2), "-")
	var s string
	for i := 0; i < len(word2); i++ {
		s = s + "/" + word2[i]
	}
	filePath := path.Join("./file/", word1[0], s, "/", word1[1])
	return true, filePath, nil, nil
}
