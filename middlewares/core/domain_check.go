package middlewares

import (
	"github.com/gin-gonic/gin"
	"GoBlade/config"
	"GoBlade/utils"
	"net/http"
	"strings"
)

func DomainGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		host := c.Request.Host
		domain := strings.Split(host, ":")[0]

		config.ConfigLock.RLock()
		allowed := config.Config.AllowedHosts
		errorDir := config.Config.ErrorDir
		config.ConfigLock.RUnlock()

		// 如果白名单为空或域名不匹配
		if len(allowed) == 0 || !isAllowedDomain(domain, allowed) {
			utils.RenderErrorPage(c, http.StatusForbidden, gin.H{
				"path": domain,
			}, errorDir)
			return
		}

		c.Next()
	}
}

func isAllowedDomain(domain string, allowed []string) bool {
	for _, item := range allowed {
		if domain == item {
			return true
		}
	}
	return false
}
