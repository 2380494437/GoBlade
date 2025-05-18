# GoBlade

🌐 轻量、清晰、可拓展的 Go 后端接口模板框架  
作者：菜鸟八哥@ChatGPT ｜ QQ：2380494437 ｜ QQ群：675644084

---

## ✨ 项目特点

- 🧱 基于 [Gin](https://github.com/gin-gonic/gin)，性能优秀
- 📁 支持多版本接口（如：/api/v1, /api/v2）
- ⚙️ 配置热重载（config.yaml）
- 🔒 内置中间件（限流、防爆、黑名单等）
- 📄 自带错误页、静态资源托管
- 🧩 支持自动建表示例（无需写 SQL）
- 💻 支持页面渲染（适合融合项目）

---

## ⚙️ 技术栈

- Go 1.20+
- Gin Web 框架
- MySQL（可选，默认关闭）
- Layui（前端页面文档）

---

## 📁 项目结构（简要）

GoBlade/
├── main.go // 启动入口
├── config/ // 配置加载、数据库连接
├── controllers/ // 控制器（按版本划分）
├── routes/ // 路由注册
├── middlewares/ // 中间件（限流、权限）
├── utils/ // 工具函数（返回结构等）
├── static/ // 前端资源（文档页面）
├── templates/ // 错误页/页面模板
├── config.yaml // 主配置文件（热重载）
└── go.mod / go.sum

---

## 🚀 快速开始

# 克隆项目
git clone https://github.com/2380494437/GoBlade.git

cd GoBlade

# 初始化依赖
go mod tidy

# 运行项目（默认监听 8080）
go run main.go
浏览器访问：
http://localhost:8080

🛠 常用命令
# 编译为可执行文件
go build -o goblade main.go

# 后台运行（Linux）
nohup ./goblade > run.log 2>&1 &

# 交叉编译（Windows → Linux）
SET GOOS=linux
SET GOARCH=amd64
go build -o goblade main.go

📚 文档入口
访问文档页面（内置 Layui 界面）实际搭建体验：
http://18.162.208.9:8080/
🔐 许可证 License
本项目使用 MIT 开源协议，商业可用，保留作者署名。

💬 社区交流
QQ：2380494437

QQ交流群：675644084

欢迎反馈问题、建议或提交 PR！
(该项目为ChatGPT完全开发：框架设计/页面/文档等全部均有ChatGPT开发)
