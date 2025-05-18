package middlewares

import (
	"github.com/gin-gonic/gin"
	"GoBlade/utils"
)

// AuthMiddleware 是一个权限中间件，开发者可自定义验证逻辑
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 示例：读取 token（你也可以换成 Session/Cookie/IP 等）
		token := c.GetHeader("Authorization")

		// TODO: 你可以替换这里的权限校验逻辑（例如：校验JWT、查询用户角色）
		if token == "" || token != "your-token" {
			utils.JSON(c, 401, "未授权的访问", nil)
			c.Abort() // 阻止继续执行后续 handler
			return
		}

		// 权限通过
		c.Next()
	}
}
