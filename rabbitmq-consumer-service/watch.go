package main

import (
	"log"

	"gopkg.in/fsnotify.v1"
)

// 监听配置文件变动
func watchConfigFile(configPath string, config *Config) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = watcher.Add(configPath)
	if err != nil {
		return err
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("配置文件已修改，正在重新加载...")
					// 配置文件变化时，重新加载
					updatedConfig, err := LoadConfig(configPath)
					if err != nil {
						log.Printf("重新加载配置失败: %s", err)
					} else {
						*config = *updatedConfig
						log.Println("配置已成功重新加载")
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("监听配置文件时出错:", err)
			}
		}
	}()
	<-done
	return nil
}
