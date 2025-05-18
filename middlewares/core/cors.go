package middlewares

import (
	"GoBlade/config"
	"github.com/gin-gonic/gin"
	"strings"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		config.ConfigLock.RLock()
		cors := config.Config.CORS
		config.ConfigLock.RUnlock()
		if !cors.Enabled {
			c.Next()
			return
		}
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			// 判断是否允许该 origin
			allowed := false
			for _, o := range cors.AllowOrigins {
				if o == "*" || o == origin {
					allowed = true
					break
				}
			}
			if allowed {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Writer.Header().Set("Vary", "Origin")
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(cors.AllowMethods, ","))
		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(cors.AllowHeaders, ","))
		c.Writer.Header().Set("Access-Control-Expose-Headers", strings.Join(cors.ExposeHeaders, ","))
		if cors.AllowCredentials {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		// 预检请求直接返回
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
