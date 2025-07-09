package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	userpb "toolbox/2025/grpc-lab/grpc-lab/proto"
)

// 实现 UserService 服务
type userServiceServer struct {
	userpb.UnimplementedUserServiceServer
}

// 实现 GetUser 方法
func (s *userServiceServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	log.Printf("Received request for user ID: %d\n", req.Id)

	// 模拟一个用户数据
	return &userpb.GetUserResponse{
		Id:    req.Id,
		Name:  "Tom",
		Email: "tom@example.com",
	}, nil
}

func main() {
	// 启动监听
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 创建 gRPC 服务
	s := grpc.NewServer()

	// 注册 UserService 服务
	userpb.RegisterUserServiceServer(s, &userServiceServer{})

	log.Println("gRPC server is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
