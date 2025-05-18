package utils

import (
	"github.com/gin-gonic/gin"
	"time"
)

// 通用 JSON 响应函数（带状态码、消息、数据、时间戳）
func JSON(c *gin.Context, code int, msg string, data any) {
	c.JSON(code, gin.H{
		"code":      code,
		"message":   msg,
		"data":      data,
		"timestamp": time.Now().Unix(), // 当前时间戳（单位：秒）
	})
}

