package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"GoBlade/config"
	"GoBlade/utils"
	"net/http"
	"time"
)

func TimeoutMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		config.ConfigLock.RLock()
		timeoutSec := config.Config.RequestTimeoutSeconds
		errorDir := config.Config.ErrorDir
		config.ConfigLock.RUnlock()

		if timeoutSec <= 0 {
			timeoutSec = 15
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(timeoutSec)*time.Second)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		done := make(chan struct{})
		panicChan := make(chan interface{}, 1)

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			c.Next()
			close(done)
		}()

		select {
		case p := <-panicChan:
			panic(p)
		case <-done:
			// 正常返回
		case <-ctx.Done():
			utils.RenderErrorPage(c, http.StatusGatewayTimeout, gin.H{
				"error": "请求处理超时",
				"path":  c.FullPath(),
			}, errorDir)
		}
	}
}
