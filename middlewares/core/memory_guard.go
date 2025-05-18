package middlewares

import (
	"GoBlade/config"
	"github.com/gin-gonic/gin"
	"runtime"
	"GoBlade/utils"
)

func MemoryGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		config.ConfigLock.RLock()
		limitMB := config.Config.MaxMemoryMB
		config.ConfigLock.RUnlock()

		// 0 表示不启用限制
		if limitMB <= 0 {
			c.Next()
			return
		}

		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		usedMB := m.Alloc / 1024 / 1024

		if int(usedMB) >= limitMB {
			
			utils.RenderErrorPage(c, 503, gin.H{
				"message": "Service Unavailable - Memory limit exceeded (%dMB used / %dMB allowed)",
			}, config.Config.ErrorDir)
			c.Abort()
			return
		}

		c.Next()
	}
}
