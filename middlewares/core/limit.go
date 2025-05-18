package middlewares

import (
	"GoBlade/config"
	"github.com/gin-gonic/gin"
	"sync"
	"GoBlade/utils"
)

// 全局信号量（懒加载）
var (
	semOnce    sync.Once
	semaphore  chan struct{}
	semLock    sync.RWMutex
	currentCap int
)

// LimitConcurrency 自动读取配置并热更新最大并发数
func LimitConcurrency() gin.HandlerFunc {
	semOnce.Do(initSemaphore)

	return func(c *gin.Context) {
		checkAndReload() // 热更新最大并发数

		semLock.RLock()
		select {
		case semaphore <- struct{}{}:
			semLock.RUnlock()
			defer func() {
				<-semaphore
			}()
			c.Next()
		default:
			semLock.RUnlock()
			utils.RenderErrorPage(c, 429, gin.H{
				"message": "Too Many Requests - Server is busy.",
			}, config.Config.ErrorDir)
			c.Abort()
		}
	}
}

// 初始化信号量
func initSemaphore() {
	cfg := config.Config
	if cfg.MaxConcurrentRequests <= 0 {
		currentCap = 0
		semaphore = make(chan struct{}, 1_000_000) // 理论上无限
	} else {
		currentCap = cfg.MaxConcurrentRequests
		semaphore = make(chan struct{}, currentCap)
	}
}

// 检查是否需要重新调整容量
func checkAndReload() {
	config.ConfigLock.RLock()
	cfg := config.Config
	config.ConfigLock.RUnlock()

	// 若配置变化，则重新创建信号量
	if cfg.MaxConcurrentRequests > 0 && cfg.MaxConcurrentRequests != currentCap {
		semLock.Lock()
		defer semLock.Unlock()
		semaphore = make(chan struct{}, cfg.MaxConcurrentRequests)
		currentCap = cfg.MaxConcurrentRequests
	}
}
