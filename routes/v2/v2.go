package v2

import (
	"github.com/gin-gonic/gin"
	// 你可以使用不同的控制器包
)

func Register(r *gin.RouterGroup) {
    // 所有业务接口必须使用 POST 方法，禁止使用 GET，以避免与静态资源 /*filepath 冲突
    // 例外情况请联系框架维护者修改静态路径映射规则
	r.POST("/install", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "v2 安装接口（示例）"})
	})
}
