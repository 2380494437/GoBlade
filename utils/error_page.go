package utils

import (
	"GoBlade/config"
	"github.com/gin-gonic/gin"
	"html/template"
	"os"
	"path/filepath"
	"fmt"
)

// RenderErrorPage 渲染错误页面
// status: HTTP 状态码（403、404、500）
// data:   模板变量（如 .error、.stack、.path）
// tmplDir: 模板目录（一般来自 config.Config.ErrorDir）
// RenderErrorPage 渲染错误页面
func RenderErrorPage(c *gin.Context, status int, data gin.H, tmplDir string) {
	// 自动查找模板名
	tmplName := getErrorTemplateName(status, tmplDir)
	tmplPath := filepath.Join(tmplDir, tmplName)

	if _, err := os.Stat(tmplPath); err == nil {
		c.Status(status)
		tmpl := template.Must(template.ParseFiles(tmplPath))
		_ = tmpl.Execute(c.Writer, data)
	} else {
		// 模板文件不存在，回退纯文本
		config.ConfigLock.RLock()
		debug := config.Config.DebugMode
		config.ConfigLock.RUnlock()

		c.Status(status)

		if status == 500 {
			if debug {
				c.String(status, "500 Internal Server Error\n%v", data)
			} else {
				c.String(status, "500 Internal Server Error - Please contact the administrator.")
			}
		} else {
			c.String(status, "%d Error - %v", status, data)
		}
	}
	c.Abort()
}
// getErrorTemplateName 返回对应错误页模板文件名
func getErrorTemplateName(code int, tmplDir string) string {
	name := fmt.Sprintf("%d.html", code)
	fullPath := filepath.Join(tmplDir, name)

	if _, err := os.Stat(fullPath); err == nil {
		return name // 例如 429.html 存在
	}
	return "error.html" // 默认兜底
}
