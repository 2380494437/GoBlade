package config

import (
	"database/sql"
	"fmt"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB // 全局共享连接池对象

func InitDB() error {
	mysql := Config.MySQL

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		mysql.User, mysql.Password, mysql.Host, mysql.Port, mysql.Dbname)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	// 设置连接池参数
	db.SetMaxOpenConns(mysql.MaxOpenConns)
	db.SetMaxIdleConns(mysql.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(mysql.ConnMaxLifetime) * time.Second)

	if err := db.Ping(); err != nil {
		return err
	}

	DB = db
	return nil
}

func ReloadDB(newCfg MySQLConfig) error {
	// 构建连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		newCfg.User, newCfg.Password, newCfg.Host, newCfg.Port, newCfg.Dbname)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	// 设置连接池参数
	db.SetMaxOpenConns(newCfg.MaxOpenConns)
	db.SetMaxIdleConns(newCfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(newCfg.ConnMaxLifetime) * time.Second)

	if err := db.Ping(); err != nil {
		return err
	}

	// 替换全局连接池
	if DB != nil {
		_ = DB.Close()
	}
	DB = db
	fmt.Println("数据库连接已重载成功")
	return nil
}
