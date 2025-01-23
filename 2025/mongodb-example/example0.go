package main

import (
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"toolbox/2025/mongodb-example/storage"
)

type User struct {
	Name   string `bson:"name"`
	Email  string `bson:"email"`
	Age    int    `bson:"age"`
	Status string `bson:"status"`
}

func main() {
	// 初始化 MongoDB 连接
	uri := "mongodb://localhost:27017"
	dbname := "test"
	mgo, err := storage.NewMongoDB(uri, dbname)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mgo.Close()

	// 批量插入文档
	/*users := []interface{}{
		User{Name: "Alice12", Email: "alice12@example.com", Age: 25, Status: "disable"},
		User{Name: "Bob21", Email: "bob21@example.com", Age: 14, Status: "disable"},
		User{Name: "Charlie44", Email: "charlie44@example.com", Age: 30, Status: "disable"},
	}
	insertManyResult, err := mgo.InsertMany("t_users", users)
	if err != nil {
		log.Fatalf("Failed to insert documents: %v", err)
	}
	fmt.Printf("Inserted documents with IDs: %v\n", insertManyResult.InsertedIDs)*/


	// 打印更新条件和操作
	fmt.Printf("Update filter: %+v\n", bson.M{"age": bson.M{"$gt": 25}})
	fmt.Printf("Update operation: %+v\n", bson.M{"$set": bson.M{"status": "active"}})

	// 批量更新文档
	updateResult, err := mgo.UpdateMany("t_users", bson.M{"age": bson.M{"$gt": 25}}, bson.M{"$set": bson.M{"status": "active"}})
	if err != nil {
		log.Fatalf("Failed to update documents: %v", err)
	}
	fmt.Printf("Updated %v documents\n", updateResult.ModifiedCount)

}
