package middlewares

import (
	"GoBlade/config"
	"github.com/gin-gonic/gin"
	"sync"
	"GoBlade/utils"
)

// IP -> 当前活跃请求数
var (
	ipCounter = make(map[string]int)
	ipLock    = sync.Mutex{}
)

func LimitConcurrencyPerIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		config.ConfigLock.RLock()
		limit := config.Config.MaxConcurrentPerIP
		config.ConfigLock.RUnlock()

		if limit <= 0 {
			c.Next()
			return
		}

		ipLock.Lock()
		count := ipCounter[ip]
		if count >= limit {
			ipLock.Unlock()
			utils.RenderErrorPage(c, 429, gin.H{
				"message": "Too Many Requests - IP limit reached",
			}, config.Config.ErrorDir)
			c.Abort()
			return
		}
		ipCounter[ip]++
		ipLock.Unlock()

		// ✅ 用 defer 确保无论什么情况都释放
		defer func() {
			ipLock.Lock()
			ipCounter[ip]--
			if ipCounter[ip] <= 0 {
				delete(ipCounter, ip)
			}
			ipLock.Unlock()
		}()

		c.Next()
	}
}

