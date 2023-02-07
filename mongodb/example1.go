package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"toolbox/conf"
	"toolbox/extend/mgo"
)

func main() {
	config := conf.GetINI()
	sec := config.Section("mongodb")
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/",
		sec.Key("username").String(),
		sec.Key("password").String(),
		sec.Key("host").String(),
		sec.Key("port").MustInt(),
	)
	dbName := "test1"
	timeout := time.Second * 30
	poolSize := 10

	collection, err := mgo.ConnMgo(uri, dbName, timeout, uint64(poolSize))
	if err != nil {
		log.Fatal(err)
	}

	c := collection.Collection("logs")

	// 插入文档
	/*logsInsert := mgo.Logs{"mongodb", "插入文档1"}
	insertResult, err := c.InsertOne(context.TODO(), logsInsert)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insertResult.InsertedID)*/

	// 批量插入文档
	/*logsManyInsert := []interface{}{
		mgo.Logs{"mongodb", "插入文档2"},
		mgo.Logs{"mongodb", "插入文档3"},
	}
	insertManyResult, err := c.InsertMany(context.TODO(), logsManyInsert)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insertManyResult.InsertedIDs)*/

	// 更新文档
	/*filter := bson.D{{"name", "mongodb"}}
	update := bson.D{{"$set", bson.D{{"data", "更新文档"}}}}
	//updateResult, err := c.UpdateOne(context.TODO(), filter, update)
	//if err != nil {
	//	log.Fatal(err)
	//}

	// 批量更新文档
	updateResult, err := c.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(updateResult.MatchedCount, updateResult.ModifiedCount)*/

	// 查询文档
	/*var logs mgo.Logs
	filter := bson.D{{"name", "mongodb"}}
	err = c.FindOne(context.TODO(), filter).Decode(&logs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(logs)*/

	// 查询所有文档
	filter := bson.D{{"name", "mongodb"}}
	cursor, err := c.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	var logsAll []mgo.Logs
	if err = cursor.All(context.TODO(), &logsAll); err != nil {
		log.Fatal(err)
	}

	for _, result := range logsAll {
		cursor.Decode(&result)
		output, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", output)
	}

}
