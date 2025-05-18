package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
	"time"
)

// 全局可热重载配置
var Config ProjectConfig
var ConfigLock sync.RWMutex

// 从 YAML 文件加载配置
func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var newConfig ProjectConfig
	if err := yaml.Unmarshal(data, &newConfig); err != nil {
		return err
	}

	// 获取旧配置用于比较
	ConfigLock.RLock()
	oldConfig := Config
	ConfigLock.RUnlock()

	// 更新配置
	ConfigLock.Lock()
	Config = newConfig
	ConfigLock.Unlock()

	fmt.Println("配置文件加载完成")

	// ==== 数据库重载逻辑 ====
	oldEnabled := oldConfig.MySQLEnabled
	newEnabled := newConfig.MySQLEnabled

	mysqlChanged := isMySQLConfigChanged(oldConfig.MySQL, newConfig.MySQL)

	if newEnabled {
		if !oldEnabled || mysqlChanged {
			fmt.Println("检测到数据库启用或配置变化，重载连接池")
			if err := ReloadDB(newConfig.MySQL); err != nil {
				fmt.Println("数据库重载失败：", err)
			}
		}
	} else {
		// 若启用状态从 true → false，关闭连接池
		if oldEnabled && DB != nil {
			_ = DB.Close()
			DB = nil
			fmt.Println("数据库已关闭（mysql_enabled=false）")
		}
	}

	return nil
}


// 开启配置热重载监听
func WatchConfig(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	err = watcher.Add(path)
	if err != nil {
		panic(err)
	}

	var reloadTimer *time.Timer

	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				// 只关注写入事件
				if ev.Op&fsnotify.Write == fsnotify.Write {
					// 重复触发则重置定时器
					if reloadTimer != nil {
						reloadTimer.Stop()
					}
					// 延迟 200ms 执行（防止多次触发）
					reloadTimer = time.AfterFunc(200*time.Millisecond, func() {
						fmt.Println("配置文件变更，自动重载中...")
						if err := LoadConfig(path); err != nil {
							fmt.Println("加载配置失败：", err)
						}
					})
				}
			case err := <-watcher.Errors:
				fmt.Println("配置监听错误:", err)
			}
		}
	}()
}


// 比较两个数据库配置是否发生变化
func isMySQLConfigChanged(oldCfg, newCfg MySQLConfig) bool {
	return oldCfg.Host != newCfg.Host ||
		oldCfg.Port != newCfg.Port ||
		oldCfg.User != newCfg.User ||
		oldCfg.Password != newCfg.Password ||
		oldCfg.Dbname != newCfg.Dbname
}
