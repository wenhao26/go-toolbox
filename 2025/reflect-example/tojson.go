// JSON 序列化
package main

import (
	"fmt"
	"reflect"
)

func toJSON(v interface{}) string {
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Struct {
		result := "{"
		for i := 0; i < value.NumField(); i++ {
			field := value.Field(i)
			if i > 0 {
				result += ","
			}
			result += fmt.Sprintf("\"%s\":\"%v\"", value.Type().Field(i).Name, field.Interface())
		}
		result += "}"

		return result
	}

	return ""
}

type User struct {
	Name string
	Age  int
}

func main() {
	user := User{
		Name: "Alice",
		Age:  20,
	}
	result := toJSON(user)
	fmt.Println(result)

	config := make(map[string]interface{})
	config["path"] = "/usr/path"
	config["size"] = 125000
	config["enable"] = true
	result1 := toJSON(config)
	fmt.Println(result1)
}
