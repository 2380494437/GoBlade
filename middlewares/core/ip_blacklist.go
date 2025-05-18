package middlewares

import (
	"bufio"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	"github.com/gin-gonic/gin"
	"GoBlade/config"
	"GoBlade/utils"
)

var (
	blacklistRules []string
	blacklistLock  sync.RWMutex
	blacklistPath  string
)

func loadBlacklist() {
	path := blacklistPath
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	var rules []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			rules = append(rules, line)
		}
	}
	blacklistLock.Lock()
	blacklistRules = rules
	blacklistLock.Unlock()
}

func isBlacklisted(ip string) bool {
	blacklistLock.RLock()
	defer blacklistLock.RUnlock()

	for _, rule := range blacklistRules {
		if strings.Contains(rule, "*") {
			if matchWildcardIP(ip, rule) {
				return true
			}
		} else if strings.Contains(rule, "/") {
			_, ipnet, err := net.ParseCIDR(rule)
			if err == nil && ipnet.Contains(net.ParseIP(ip)) {
				return true
			}
		} else if rule == ip {
			return true
		}
	}
	return false
}

func matchWildcardIP(ip, pattern string) bool {
	ipParts := strings.Split(ip, ".")
	patternParts := strings.Split(pattern, ".")
	if len(ipParts) != 4 || len(patternParts) != 4 {
		return false
	}
	for i := 0; i < 4; i++ {
		if patternParts[i] != "*" && patternParts[i] != ipParts[i] {
			return false
		}
	}
	return true
}

// 黑名单 IP 中间件（自动从 config 读取路径）
func IPBlacklistMiddleware() gin.HandlerFunc {
	// 获取路径
	config.ConfigLock.RLock()
	blacklistPath = config.Config.IPBlacklistFile
	config.ConfigLock.RUnlock()

	// 初始加载
	loadBlacklist()

	// 后台定时热加载
	go func() {
		for {
			time.Sleep(time.Minute)
			loadBlacklist()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if isBlacklisted(ip) {
			utils.RenderErrorPage(c, http.StatusForbidden, gin.H{
				"path": "黑名单 IP：" + ip,
			}, config.Config.ErrorDir)
			c.Abort()
			return
		}
		c.Next()
	}
}
