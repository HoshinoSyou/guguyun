package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"
)

type Header struct {
	Typ string `json:"typ"`
	Alg string `json:"alg"`
}

type Payload struct {
	Iss      string `json:"iss"`
	Exp      string `json:"exp"`
	Iat      string `json:"iat"`
	Username string `json:"username"`
}

//CreateHeader 创建Token头部
func CreateHeader() string {
	h := Header{
		Typ: "GROM",
		Alg: "HS256",
	}
	bytes, err := json.Marshal(h)
	if err != nil {
		log.Printf("marshal failed!序列化失败！错误信息：%v", err)
		return ""
	}
	header := base64.StdEncoding.EncodeToString(bytes)
	return header
}

//CreatePayload 创建Token载荷
func CreatePayload(username string) string {
	p := Payload{
		Iss:      "HoshinoSyou",
		Exp:      strconv.FormatInt(time.Now().Add(15*24*time.Hour).Unix(), 10),
		Iat:      strconv.FormatInt(time.Now().Unix(), 10),
		Username: username,
	}
	bytes, err := json.Marshal(p)
	if err != nil {
		log.Printf("marshal failed!序列化失败！错误信息：%v", err)
		return ""
	}
	payload := base64.StdEncoding.EncodeToString(bytes)
	return payload
}

//CreateSignature 创建签证
func CreateSignature(username string) string {
	h := CreateHeader()
	p := CreatePayload(username)
	str := strings.Join([]string{h, p}, ".")
	key := "HoshinoSyou"
	hash := hmac.New(sha256.New, []byte(key))
	hash.Write([]byte(str))
	s := hash.Sum(nil)
	signature := base64.StdEncoding.EncodeToString(s)
	return signature
}

//Create 登陆后创建Token
func Create(username string) (token string) {
	header := CreateHeader()
	payload := CreatePayload(username)
	signature := CreateSignature(username)
	token = strings.Join([]string{header, payload}, ".") + "." + signature
	return token
}
