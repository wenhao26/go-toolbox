package main

import (
	"github.com/tidwall/gjson"
)

// GJSON是一个Go包，它提供了一种快速简单的方法来从json文档中获取值
// 它具有单行检索、点符号路径、迭代和解析json行等功能
func main() {
	json := `{"code":0,"message":"success","data":[{"category_id":24,"category_name":"Werewolf/Vampire","category_cover":""},{"category_id":23,"category_name":"Romance","category_cover":""},{"category_id":26,"category_name":"Action/Adventure","category_cover":""},{"category_id":34,"category_name":"Science Fiction","category_cover":""},{"category_id":35,"category_name":"Humor","category_cover":""},{"category_id":36,"category_name":"Historical","category_cover":""},{"category_id":37,"category_name":"General Fiction","category_cover":""},{"category_id":38,"category_name":"Adult Romance","category_cover":""},{"category_id":39,"category_name":"New Adult","category_cover":""},{"category_id":40,"category_name":"LGBTQ+","category_cover":""}],"_lang":"en","_time":1682578404}`

	value := gjson.Get(json, "data")
	println(value.String())
}
