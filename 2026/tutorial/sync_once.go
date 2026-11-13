package main

import (
	"fmt"
	"sync"
)

type DBConnection struct {
	DSN string
}

var (
	dbInstance *DBConnection
	once       sync.Once
)

// GetDBInstance 获取全局唯一的数据库连接实例
func GetDBInstance() *DBConnection {
	// once.Do 保证传入的闭包函数在整个程序生命周期中只执行一次
	// 它是并发安全的，性能极高
	once.Do(func() {
		fmt.Println("初始化数据库连接...")
		dbInstance = &DBConnection{
			DSN: "mysql://root:pass@localhost:3306/go",
		}
	})
	return dbInstance
}

// - 线程安全的懒加载单例与初始化
// 场景：在企业级应用中，有些全局资源（如配置中心连接、加密密匙、全局缓存客户端）只需要初始化一次，
// 且希望在第一次真正用到时才去加载（Lazy Load），同时要保证在高并发下的绝对安全。
//
// 惊艳之处：标准库的 sync.Once 内部使用了双重检查锁（Double-Checked Locking）的思想，但代码极其简短优雅，完全不需要我们手写复杂的互斥锁逻辑。
//
// 学习收获：学会用 sync.Once 代替粗暴的全局锁，既保证了线程安全，又避免了每次访问都要加锁带来的性能损耗。

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			db := GetDBInstance()
			fmt.Printf("协程 %d 获取到的实例地址：%p\n", i, db)
		}()
	}
	wg.Wait()

	fmt.Println("结束")
}
