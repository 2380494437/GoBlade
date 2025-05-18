package middlewares

import (
	"GoBlade/config"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"GoBlade/utils"
)

func BodyLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		config.ConfigLock.RLock()
		limit := config.Config.MaxBodySize
		config.ConfigLock.RUnlock()

		if limit <= 0 {
			c.Next()
			return
		}

		// 如果 Content-Length 存在且超限，提前拒绝
		if cl := c.Request.Header.Get("Content-Length"); cl != "" {
			if length, err := strconv.ParseInt(cl, 10, 64); err == nil && length > limit {
				//c.String(http.StatusRequestEntityTooLarge, "413 Request Entity Too Large")
				utils.RenderErrorPage(c, 413, gin.H{
					"message": "Request Entity Too Large",
				}, config.Config.ErrorDir)
				c.Abort()
				return
			}
		}

		// 强制限制实际读取的最大体积（即使 Content-Length 被伪造）
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)

		// 手动触发读取以便触发 MaxBytesReader 的限制（可选）
		if c.Request.ContentLength > 0 {
			_, err := io.ReadAll(io.LimitReader(c.Request.Body, 1))
			if err != nil {
				//c.String(http.StatusRequestEntityTooLarge, "413 Body too large")
				utils.RenderErrorPage(c, 413, gin.H{
					"message": "Request Entity Too Large",
				}, config.Config.ErrorDir)
				c.Abort()
				return
			}
			// 重置读取流位置
			c.Request.Body.Close()
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)
		}

		c.Next()
	}
}
