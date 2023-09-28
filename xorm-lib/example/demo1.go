package main

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	_ "xorm.io/xorm/caches"
	"xorm.io/xorm/names"
)

type Article struct {
	Id           int64  `xorm:"'id'"`
	CreatedAt    string `xorm:"'created_at'"`
	UpdatedAt    string `xorm:"'updated_at'"`
	DeletedAt    string `xorm:"'deleted_at'"`
	Title        string `xorm:"'title'"`
	Cid          int64  `xorm:"'cid'"`
	Desc         string `xorm:"'desc'"`
	Content      string `xorm:"'content'"`
	Img          string `xorm:"'img'"`
	CommentCount int64  `xorm:"'comment_count'"`
	ReadCount    int64  `xorm:"'read_count'"`
}

var engine *xorm.Engine

func main() {
	t := time.Now()

	// 创建引擎
	var err error
	engine, err = xorm.NewEngine("mysql", "root:root@/testdb?charset=utf8")
	if err != nil {
		panic(err)
	}
	defer engine.Close()

	engine.ShowSQL(true)

	//engine.SetMapper(names.GonicMapper{})
	// 表前缀
	tbMapper := names.NewPrefixMapper(names.SnakeMapper{}, "t_")
	engine.SetTableMapper(tbMapper)

	// 开启缓存
	//cacher := caches.NewLRUCacher(caches.NewMemoryStore(), 1000)
	//engine.SetDefaultCacher(cacher)
	//engine.MapCacher(&Article{}, cacher)

	// 查询
	//var article Article
	//engine.Where("id=?", 1).Get(&article)
	//engine.SQL("select * from t_article where id=1").Get(&article)

	var article []Article
	engine.Limit(1000).Find(&article)

	for _, v := range article {
		fmt.Printf("id=%d,title=%s,create-at=%s \n", v.Id, v.Title, v.CreatedAt)
	}

	fmt.Println(time.Since(t))
}
