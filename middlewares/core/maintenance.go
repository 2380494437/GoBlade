package middlewares

import (
	"GoBlade/config"
	"GoBlade/utils"
	"github.com/gin-gonic/gin"
)

func MaintenanceMode() gin.HandlerFunc {
	return func(c *gin.Context) {
		config.ConfigLock.RLock()
		enabled := config.Config.MaintenanceMode
		errorDir := config.Config.ErrorDir
		config.ConfigLock.RUnlock()

		if enabled {
			utils.RenderErrorPage(c, 521, gin.H{
				"error": "系统维护中，请稍后再访问。",
			}, errorDir)
			return
		}

		c.Next()
	}
}
