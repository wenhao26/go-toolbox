package mgo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 连接池连接方式
func ConnMgo(uri, dbName string, timeout time.Duration, poolSize uint64) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	clientOptions.SetMaxPoolSize(poolSize)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	return client.Database(dbName), nil
}

// 演示例子
/*func NewMgo() {
	config := conf.GetINI()
	sec := config.Section("mongodb")

	// mongodb://用户:密码@地址:27017/数据库?connect=direct
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/",
		sec.Key("username").String(),
		sec.Key("password").String(),
		sec.Key("host").String(),
		sec.Key("port").MustInt(),
	)
	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI(uri)

	// 连接mongodb
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// 连接检查
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("连接mongodb成功！")

	// 指定获取要操作的数据集
	collection := client.Database("test1").Collection("t_logs")
	fmt.Println(collection)

	// 断开连接
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}*/
