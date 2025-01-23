// Package mongodb 提供了一个简单且高效的 MongoDB 操作封装。
// 它支持常见的 CRUD 操作（插入、查询、更新、删除），并提供了批量插入和批量更新的功能。
// 该包还封装了连接管理、超时控制和错误处理，适合在生产环境中使用。
//
// 主要功能：
//   - 连接 MongoDB 数据库。
//   - 插入单条或多条文档。
//   - 查询单条或多条文档。
//   - 更新单条或多条文档。
//   - 删除单条文档。
//   - 统计文档数量。
//
// 示例用法：
//   // 初始化 MongoDB 连接
//   mongoDB, err := mongodb.NewMongoDB("mongodb://localhost:27017", "testdb")
//   if err != nil {
//       log.Fatalf("Failed to connect to MongoDB: %v", err)
//   }
//   defer mongoDB.Close()
//
//   // 插入一条文档
//   result, err := mongoDB.InsertOne("users", bson.M{"name": "Alice", "age": 25})
//   if err != nil {
//       log.Fatalf("Failed to insert document: %v", err)
//   }
//   fmt.Printf("Inserted document with ID: %v\n", result.InsertedID)
//
// 注意事项：
//   - 确保 MongoDB 服务已启动并可以访问。
//   - 在生产环境中，建议为关键操作设置合理的超时时间。
//   - 批量操作时，注意控制每次操作的数据量，避免内存占用过高或操作超时。
//
// 依赖：
//   - Go MongoDB 驱动：go.mongodb.org/mongo-driver/mongo
//
// 作者：Wenhao Wu
// 版本：v1.0.0
package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB mgo 存储结构体
type MongoDB struct {
	client   *mongo.Client
	database *mongo.Database
}

// NewMongoDB 创建 mgo 实例
func NewMongoDB(uri, dbname string) (*MongoDB, error) {
	// 设置连接超时
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 连接MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// 检查连接是否成功
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB successfully!")

	return &MongoDB{
		client:   client,
		database: client.Database(dbname),
	}, nil
}

// Close 关闭 mgo 连接
func (m *MongoDB) Close() {
	if m.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = m.client.Disconnect(ctx)
		log.Println("MongoDB connection closed.")
	}
}

// getCollectionWithContext 获取集合并创建上下文
// collectionName: 集合名称
func (m *MongoDB) getCollectionWithContext(collectionName string) (*mongo.Collection, context.Context, context.CancelFunc) {
	collection := m.database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()

	return collection, ctx, cancel
}

// InsertOne 插入一条文档数据
// collectionName: 集合名称
// document: 要插入的文档（可以是结构体或者 bson.M）
func (m *MongoDB) InsertOne(collectionName string, document interface{}) (*mongo.InsertOneResult, error) {
	//collection := m.database.Collection(collectionName)
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	collection, ctx, cancel := m.getCollectionWithContext(collectionName)
	defer cancel()

	result, err := collection.InsertOne(ctx, document)
	if err != nil {
		return nil, fmt.Errorf("failed to insert document: %v", err)
	}

	return result, nil
}

// InsertMany 批量插入文档数据
// collectionName: 集合名称
// document: 要插入的文档列表（[]interface{}）
func (m *MongoDB) InsertMany(collectionName string, documents []interface{}) (*mongo.InsertManyResult, error) {
	collection := m.database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.InsertMany(ctx, documents)
	if err != nil {
		return nil, fmt.Errorf("failed to insert documents: %v", err)
	}

	return result, nil
}

// FindOne 查询一条文档
// collectionName: 集合名称
// filter: 查询条件（bson.M）
// result: 用于存储查询结果的指针
func (m *MongoDB) FindOne(collectionName string, filter bson.M, result interface{}) error {
	collection := m.database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, filter).Decode(result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("document not found")
		}
		return fmt.Errorf("failed to find document: %v", err)
	}

	return nil
}

// FindMany 查询多条文档
// collectionName: 集合名称
// filter: 查询条件（bson.M）
// result: 用于存储查询结果的指针
func (m *MongoDB) FindMany(collectionName string, filter bson.M, results interface{}) error {
	collection := m.database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to find documents: %v", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, results); err != nil {
		return fmt.Errorf("failed to decode documents: %v", err)
	}

	return nil
}

// UpdateOne 更新一条文档
// collectionName: 集合名称
// filter: 查询条件（bson.M）
// update: 更新操作（bson.M）
func (m *MongoDB) UpdateOne(collectionName string, filter, update bson.M) (*mongo.UpdateResult, error) {
	collection := m.database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return nil, fmt.Errorf("failed to update document: %v", err)
	}

	return result, nil
}

// UpdateMany 批量更新文档
// collectionName: 集合名称
// filter: 查询条件（bson.M）
// update: 更新操作（bson.M）
func (m *MongoDB) UpdateMany(collectionName string, filter, update bson.M) (*mongo.UpdateResult, error) {
	collection := m.database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.UpdateMany(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return nil, fmt.Errorf("failed to update documents: %v", err)
	}

	return result, nil
}

// DeleteOne 删除一条文档
// collectionName: 集合名称
// filter: 查询条件（bson.M）
func (m *MongoDB) DeleteOne(collectionName string, filter bson.M) (*mongo.DeleteResult, error) {
	collection := m.database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to delete document: %v", err)
	}

	return result, nil
}

// CountDocuments 统计文档数量
// collectionName: 集合名称
// filter: 查询条件（bson.M）
func (m *MongoDB) CountDocuments(collectionName string, filter bson.M) (int64, error) {
	collection := m.database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count documents: %v", err)
	}

	return count, nil
}
