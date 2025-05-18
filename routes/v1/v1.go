package v1

import (
	"github.com/gin-gonic/gin"
	"GoBlade/controllers/v1"
	"GoBlade/middlewares"
	"GoBlade/utils"
	"GoBlade/config"
	"time"
)

// 示例动态获取热重载配置,用于获取数据表表名
func getCfg() config.ProjectConfig {
	config.ConfigLock.RLock()
	defer config.ConfigLock.RUnlock()
	return config.Config
}

// 注册 v1 路由
func Register(r *gin.RouterGroup) {
	// 安装接口（自动建表）
	r.POST("/install", controllers.InstallHandler)

	// 崩溃测试接口
	r.GET("/test-panic", func(c *gin.Context) {
		panic("测试崩溃：这是一条人为触发的 panic")
	})

	// 用户相关接口（需鉴权）
	protected := r.Group("/user")
	protected.Use(middlewares.AuthMiddleware())
	{
		//示例一个路由方法 对应 /api/v1/user/create
		protected.POST("/create", func(c *gin.Context) {
			utils.JSON(c, 200, "操作成功", gin.H{"id": 1})
		})
		
		//示例一个路由方法 对应 /api/v1/user/list
		protected.POST("/list", func(c *gin.Context) {
			time.Sleep(5 * time.Second) // 阻塞模拟耗时
			cfg := getCfg()//获取重载后的配置文件
			db := config.DB//从连接池里拿一个数据库句柄
			table := cfg.TablePrefix + "user"//拼接表名
		
			// 假设前端传了关键词参数 keyword
			keyword := c.PostForm("keyword")
			if keyword == "" {
				keyword = "%" // 查询全部
			} else {
				keyword = "%" + keyword + "%"
			}
		
			query := "SELECT id, nickname, account FROM " + table + " WHERE account LIKE ?"
			rows, err := db.Query(query, keyword)
			if err != nil {
				utils.JSON(c, 500, "数据库查询失败："+err.Error(), nil)
				return
			}
			defer rows.Close()
		
			var users []gin.H
			for rows.Next() {
				var id int
				var nickname, account string
				if err := rows.Scan(&id, &nickname, &account); err != nil {
					utils.JSON(c, 500, "数据解析失败："+err.Error(), nil)
					return
				}
				users = append(users, gin.H{
					"id":       id,
					"nickname": nickname,
					"account":  account,
				})
			}
		
			utils.JSON(c, 200, "查询成功", users)
		})
		
	}
}
