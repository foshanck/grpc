package main

import (
	"context"
	"github.com/ckLearning/grpc/deviceinfo"
	"github.com/golang/protobuf/ptypes/empty"
	"io"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 定义服务端使用的 Unix Socket 路径[6,7](@ref)
	socketPath := "/tmp/grpc_uds/deviceinfo.sock"

	// 创建自定义的拨号器
	dialOption := grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
		// 注意：这里的 addr 参数通常由 grpc.Dial 的 target 参数提供，但使用 UDS 时我们更直接使用 socketPath
		// 设置连接超时
		timeout := 30 * time.Second
		return net.DialTimeout("unix", socketPath, timeout)
	})

	// 建立连接
	// target 参数使用 "unix://" 方案或任意占位符，因为实际连接由自定义拨号器处理
	conn, err := grpc.Dial("unix:///tmp/grpc_uds/deviceinfo.sock",
		dialOption,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // 禁用 TLS
		// 可根据需要添加其他 DialOption，如负载均衡、拦截器等
	)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	// 创建 gRPC 客户端实例[2,3](@ref)
	client := deviceinfo.NewDeviceInfoClient(conn)

	// 调用 Subscribe 方法，获取消息流[5,8](@ref)
	/*podDeviceInfos, err := client.ListAllMigDeviceInfos(context.Background(),  &empty.Empty{

	})*/


	stream, err := client.WatchMigDeviceInfos(context.Background(), &empty.Empty{
	})
	if err != nil {
		log.Fatalf("watch mig devices failed. Err: %+v", err)
	}
	for {
		podDeviceInfos, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.List all mig deviices failed(_) = _, %v", client, err)
		}
		print(podDeviceInfos)
	}

	// 等待一段时间，模拟客户端持续运行
	time.Sleep(1 * time.Second)
	log.Println("Client exiting.")
}

func print(podDeviceInfos *deviceinfo.PodDeviceInfos){
	if podDeviceInfos != nil && len(podDeviceInfos.PodDeviceinfos) > 0 {
		for podName,podDeviceInfo := range podDeviceInfos.PodDeviceinfos {
			log.Printf("Pod [%s]", podName)
			for containerName, cds := range podDeviceInfo.ContainerDevices {
				for _, containerDevice := range cds.ContainerDevices {
					log.Printf("Container [%s], gpu id [%s], mig id [%s]", containerName, containerDevice.ParetDeviceId, containerDevice.MigDeviceId)
				}
			}
		}
	}

}