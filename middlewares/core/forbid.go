package middlewares

import (
	"github.com/gin-gonic/gin"
	"GoBlade/config"
	"GoBlade/utils"
	"net/http"
	"strings"
)

// 中间件：拦截敏感路径访问（如 .git、.env 等）
func ForbidSensitivePaths() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		config.ConfigLock.RLock()
		cfg := config.Config
		errorDir := cfg.ErrorDir
		config.ConfigLock.RUnlock()

		// 1. 匹配完整路径
		for _, p := range cfg.ForbiddenPaths {
			if path == p {
				utils.RenderErrorPage(c, http.StatusForbidden, gin.H{
					"path": path,
				}, errorDir)
				return
			}
		}

		// 2. 匹配后缀
		for _, suffix := range cfg.ForbiddenSuffixes {
			if strings.HasSuffix(path, suffix) {
				utils.RenderErrorPage(c, http.StatusForbidden, gin.H{
					"path": path,
				}, errorDir)
				return
			}
		}

		// 3. 匹配目录
		for _, dir := range cfg.ForbiddenDirs {
			if strings.HasPrefix(path, dir+"/") || path == dir {
				utils.RenderErrorPage(c, http.StatusForbidden, gin.H{
					"path": path,
				}, errorDir)
				return
			}
		}

		c.Next()
	}
}
