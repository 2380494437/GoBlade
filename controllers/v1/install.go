package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"GoBlade/config"
	"GoBlade/utils"
)

type TableDefinition struct {
	Remark string
	Fields map[string]string // 字段名 -> 字段完整SQL定义（不含字段名）
}

func InstallHandler(c *gin.Context) {
	config.ConfigLock.RLock()
	cfg := config.Config
	config.ConfigLock.RUnlock()
	db := config.DB

	if !cfg.MySQLEnabled || db == nil {
		utils.JSON(c, 500, "数据库未启用或连接失败", nil)
		return
	}
    
	// 通用建表代码,只需要在这个数组中定义表结构即可==========================
	tableDefinitions := map[string]TableDefinition{
    	"user": {
    		Remark: "用户表",
    		Fields: map[string]string{
    			"id":       "INT AUTO_INCREMENT PRIMARY KEY COMMENT '自增索引id'",
    			"nickname": "VARCHAR(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '昵称'",
    			"account":  "VARCHAR(20) NOT NULL UNIQUE COMMENT '账号(仅数字，最长20位)'",
    			"password": "VARCHAR(255) NOT NULL COMMENT '密码(加密后长字符串)'",
    		},
    	},
    	"visit_log": {
    		Remark: "访问日志表",
    		Fields: map[string]string{
    			"id":         "INT AUTO_INCREMENT PRIMARY KEY COMMENT '访问日志id'",
    			"ip":         "VARCHAR(50) NOT NULL COMMENT 'IP地址'",
    			"endpoint":   "VARCHAR(255) NOT NULL COMMENT '请求路径'",
    			"user_agent": "VARCHAR(255) COMMENT 'User-Agent'",
    			"created_at": "DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '访问时间'",
    		},
    	},
    }
    //============================================================下方代码勿动

	for rawTableName, def := range tableDefinitions {
		tableName := cfg.TablePrefix + rawTableName

		// 检查表是否存在
		var tmp string
		query := fmt.Sprintf("SHOW TABLES LIKE '%s'", tableName)
		err := db.QueryRow(query).Scan(&tmp)

		tableExists := err == nil

		if !tableExists {
			// 构建建表字段语句
			fieldsSQL := ""
            for name, fieldDef := range def.Fields {
            	fieldsSQL += fmt.Sprintf("`%s` %s,\n", name, fieldDef)
            }
			fieldsSQL = fieldsSQL[:len(fieldsSQL)-2] // 去掉最后一个逗号

			createSQL := fmt.Sprintf(`
				CREATE TABLE %s (
					%s
				) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='%s';`, tableName, fieldsSQL, def.Remark)

			if _, err := db.Exec(createSQL); err != nil {
				c.JSON(500, gin.H{"code": 500, "message": fmt.Sprintf("创建表 %s 失败：%s", tableName, err.Error())})
				return
			}
			fmt.Println("✅ 创建表成功：", tableName)
		} else {
			// 检查每个字段是否存在，自动补全
			for fieldName, fieldDef := range def.Fields {
				checkSQL := fmt.Sprintf("SHOW COLUMNS FROM `%s` LIKE '%s'", tableName, fieldName)
                rows, err := db.Query(checkSQL)
				if err != nil {
					c.JSON(500, gin.H{"code": 500, "message": fmt.Sprintf("查询字段 %s 错误：%s", fieldName, err.Error())})
					return
				}
				defer rows.Close()

				if !rows.Next() {
					// 字段不存在，添加
					addSQL := fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `%s` %s;", tableName, fieldName, fieldDef)
					if _, err := db.Exec(addSQL); err != nil {
						c.JSON(500, gin.H{"code": 500, "message": fmt.Sprintf("添加字段 %s 失败：%s", fieldName, err.Error())})
						return
					}
					fmt.Printf("✅ 已为表 %s 添加缺失字段：%s\n", tableName, fieldName)
				}
			}
		}
	}

	c.JSON(200, gin.H{"code": 200, "message": "表结构已创建或补全成功"})
}
