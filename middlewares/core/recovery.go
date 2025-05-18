package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"GoBlade/config"
	"GoBlade/utils"
	"net/http"
	"runtime/debug"
)

func CustomRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())

				// 读取配置（锁定）
				config.ConfigLock.RLock()
				cfg := config.Config
				config.ConfigLock.RUnlock()

				// 控制台输出（仅调试模式）
				if cfg.DebugMode {
					fmt.Printf("[Panic Recovered] %v\n%s\n", err, stack)
				}

				// 使用统一的错误页面渲染
				utils.RenderErrorPage(c, http.StatusInternalServerError, gin.H{
					"error": err,
					"stack": stack,
					"debug": cfg.DebugMode,
				}, cfg.ErrorDir)
			}
		}()
		c.Next()
	}
}
