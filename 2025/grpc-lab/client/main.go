package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"

	userpb "toolbox/2025/grpc-lab/grpc-lab/proto"
)

func main() {
	// 建立连接（非 TLS）
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// 创建 gRPC 客户端
	client := userpb.NewUserServiceClient(conn)

	// 设置请求超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 发起调用
	resp, err := client.GetUser(ctx, &userpb.GetUserRequest{Id: 123})
	if err != nil {
		log.Fatalf("could not get user: %v", err)
	}

	log.Printf("User Info: ID=%d, Name=%s, Email=%s", resp.Id, resp.Name, resp.Email)
}