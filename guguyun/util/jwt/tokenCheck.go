package jwt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"guguyun/util/response"
	"log"
	"strings"
)

//CheckToken 核对检查token头部、载荷、签证三部分信息
func CheckToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("Authorization")
		split := strings.Split(token, ".")
		if len(split) != 3 {
			err := errors.New("token构建错误")
			log.Println(err)
			response.Error(ctx, "token构建错误",err)
			return
		}
		_, err := base64.StdEncoding.DecodeString(split[0])
		if err != nil {
			err = errors.New("header解析错误")
			log.Println(err)
			response.Error(ctx, "header解析错误",err)
			ctx.Abort()
			return
		}
		p, err := base64.StdEncoding.DecodeString(split[1])
		if err != nil {
			err = errors.New("payload解析错误")
			log.Println(err)
			response.Error(ctx, "payload解析错误",err)
			ctx.Abort()
			return
		}
		_, err = base64.StdEncoding.DecodeString(split[2])
		if err != nil {
			err = errors.New("signature解析错误")
			log.Println(err)
			response.Error(ctx, "signature解析错误",err)
			ctx.Abort()
			return
		}
		var payload Payload
		err = json.Unmarshal(p, &payload)
		if err != nil {
			log.Printf("unmarshal failed!反序列化失败！错误信息：%v", err)
			response.Error(ctx, "反序列化失败",err)
			ctx.Abort()
			return
		}
		ctx.Set("username", payload.Username)
		ctx.Next()
	}
}
