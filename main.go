package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"GoBlade/config"
	"GoBlade/routes"
	"GoBlade/middlewares/core"
	"net/http"
	"os"
	"runtime"
	"path/filepath"
	"GoBlade/utils"
)

func main() {
	
	// === 启动信息 Banner ===
	_ = config.LoadConfig("config.yaml")         // 第一次加载
	config.WatchConfig("config.yaml")            // 启动热重载
	fmt.Println("GoBlade Web 框架 启动完成...")
	fmt.Println("\n\n\n")
	cfg := config.Config
	fmt.Println("┌──────────────────────────────────────────────────────┐")
	fmt.Println("│             GoBlade - 一个轻量级 Web 框架            │")
	fmt.Println("├──────────────────────────────────────────────────────┤")
	fmt.Printf("│ 项目名称       : %-36s│\n", cfg.ProjectName)
	fmt.Printf("│ 监听地址       : %-36s│\n", fmt.Sprintf("%s:%d", cfg.BindIP, cfg.Port))
	fmt.Printf("│ 调试模式       : %-36v│\n", cfg.DebugMode)
	fmt.Printf("│ 静态目录       : %-36s│\n", cfg.StaticDir)
	fmt.Printf("│ 错误页面       : %-36s│\n", cfg.ErrorDir+"/")
	fmt.Printf("│ 数据库         : %-36v│\n", cfg.MySQLEnabled)
	fmt.Printf("│ 系统版本       : %-36s│\n", runtime.GOOS+" / "+runtime.GOARCH)
	fmt.Printf("│ 环境版本       : %-36s│\n", runtime.Version())
	fmt.Println("├──────────────────────────────────────────────────────┤")
	//fmt.Println("└──────────────────────────────────────────────────────┘")

	if cfg.DebugMode {
		fmt.Println("⚠️  当前为开发模式，建议线上使用 ReleaseMode")
	}

	// 根据配置设置 Gin 的运行模式
	if cfg.DebugMode {
		// 开启调试模式，输出详细日志
		gin.SetMode(gin.DebugMode)
	} else {
		// 生产模式，不输出详细日志
		gin.SetMode(gin.ReleaseMode)
	}
	// 初始化 Gin 服务
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(middlewares.CustomRecovery())//错误捕获
	r.Use(middlewares.CORSMiddleware())//跨域
	r.Use(middlewares.DomainGuard())//域名拦截
	r.Use(middlewares.IPBlacklistMiddleware())//黑名单IP拦截
	r.Use(middlewares.LimitConcurrency())//总并发限制
	r.Use(middlewares.LimitConcurrencyPerIP())//单IP并发限制
	r.Use(middlewares.BodyLimit())//最大请求体限制
	r.Use(middlewares.MemoryGuard())//最大内存限制
	r.Use(middlewares.ForbidSensitivePaths())//路由规则拦截
	r.Use(middlewares.TimeoutMiddleware())//请求超时规则拦截
	r.Use(middlewares.MaintenanceMode())//维护模式
	r.Use(middlewares.GzipMiddleware())//gzip压缩

	routes.RegisterRoutes(r)
	r.NoRoute(func(c *gin.Context) {
		cfg := config.Config
		reqPath := filepath.Clean(c.Request.URL.Path)
	
		// 1. 构造完整物理路径（静态目录 + 请求路径）
		fullPath := filepath.Join(cfg.StaticDir, reqPath)
	
		// 2. 若访问路径是目录，尝试拼接默认文档
		if stat, err := os.Stat(fullPath); err == nil && stat.IsDir() {
			fullPath = filepath.Join(fullPath, cfg.DefaultDocument)
		}
	
		// 3. 若文件存在，返回文件内容
		if stat, err := os.Stat(fullPath); err == nil && !stat.IsDir() {
			c.File(fullPath)
			return
		}
	
		// 4. 否则返回 404 错误页
		utils.RenderErrorPage(c, http.StatusNotFound, gin.H{
			"path": c.Request.URL.Path,
		}, cfg.ErrorDir)
		
		
	})
	
	

	
	addr := fmt.Sprintf("%s:%d", cfg.BindIP, cfg.Port)
	fmt.Printf("│ Gin Server     : %-36s│\n", "OK")
	fmt.Println("├──────────────────────────────────────────────────────┤")
	fmt.Printf("│ 作者           : %-32s│\n", "菜鸟八哥")
	fmt.Printf("│ QQ             : %-36s│\n", "2380494437")
	fmt.Printf("│ QQ群           : %-36s│\n", "675644084")
	fmt.Println("└──────────────────────────────────────────────────────┘")
	//r.Run(addr)
	if err := r.Run(addr); err != nil {
		fmt.Println("\n\n\n")
		fmt.Printf("Gin服务启动失败,原因： %s: %v\n", addr, err)
		os.Exit(1)
	}
}
