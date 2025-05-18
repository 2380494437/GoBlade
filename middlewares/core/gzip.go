package middlewares

import (
	"compress/gzip"
	"GoBlade/config"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

func GzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		config.ConfigLock.RLock()
		enabled := config.Config.GzipEnabled
		config.ConfigLock.RUnlock()

		// 若未启用或客户端不支持 gzip，则跳过
		if !enabled || !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		// 设置头部
		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")

		// 替换 writer
		gz := gzip.NewWriter(c.Writer)
		defer gz.Close()

		c.Writer = &gzipWriter{
			ResponseWriter: c.Writer,
			Writer:         gz,
		}

		c.Next()
	}
}

// 自定义 writer：压缩响应体
type gzipWriter struct {
	gin.ResponseWriter
	io.Writer
}

func (w *gzipWriter) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
}
