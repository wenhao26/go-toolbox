package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func main() {
	// 监控对象
	paths := []string{
		"F:\\",
	}

	// 创建监控对象
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln(err)
	}
	defer watcher.Close()

	// 添加需要监控的对象
	for _, path := range paths {
		log.Printf("监听对象：%s\n", path)

		err := watcher.Add(path)
		if err != nil {
			log.Fatalln(err)
		}
	}

	// 另启一个协程来处理监控对象事件
	go func() {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()
		for {
			select {
			case event := <-watcher.Events:
				{
					if event.Op&fsnotify.Create == fsnotify.Create {
						log.Println("创建文件：", event.Name)
					}
					if event.Op&fsnotify.Write == fsnotify.Write {
						log.Println("更新文件内容：", event.Name)
					}
					if event.Op&fsnotify.Remove == fsnotify.Remove {
						log.Println("删除文件：", event.Name)
					}
					if event.Op&fsnotify.Rename == fsnotify.Rename {
						log.Println("文件重命名：", event.Name)
					}
					if event.Op&fsnotify.Chmod == fsnotify.Chmod {
						log.Println("修改文件权限：", event.Name)
					}
				}
			case err := <-watcher.Errors:
				{
					log.Fatalln("异常情况：", err)
					return
				}
			}
		}
	}()

	select {}
}
