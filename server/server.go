package main

import (

	"github.com/ckLearning/grpc/deviceinfo"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"path/filepath"

)

func main() {
	startServer()
}

func startServer() {

	socketDir := "/tmp/grpc_uds"
	socketPath := filepath.Join(socketDir, "deviceinfo.sock")

	// 确保 socket 目录存在
	if err := os.MkdirAll(socketDir, 0755); err != nil {
		log.Fatalf("Failed to create socket directory: %v", err)
	}

	// 清理可能存在的旧 socket 文件[6,7](@ref)
	if err := os.RemoveAll(socketPath); err != nil {
		log.Fatalf("Failed to clean up old socket: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	deviceinfo.RegisterDeviceInfoServer(grpcServer, newServer())

	// 在 Unix Socket 上监听[6,7](@ref)
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatalf("Failed to listen on Unix Socket: %v", err)
	}
	defer listener.Close()

	// 设置 Socket 文件权限（仅限当前用户读写）[6,7](@ref)
	if err := os.Chmod(socketPath, 0600); err != nil {
		log.Fatalf("Failed to set socket file permissions: %v", err)
	}

	log.Printf("gRPC server listening on unix socket: %s", socketPath)

	// 启动 gRPC 服务器[2,3](@ref)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}

func newServer() deviceinfo.DeviceInfoServer{

	return &DevicePluginServer{}
}
