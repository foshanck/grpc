package main

import (
	"context"
	"github.com/ckLearning/grpc/deviceinfo"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"strconv"
	"sync"
	"time"
)

//var _ deviceinfo.DeviceInfoServer = DevicePluginServer()
var (
	requestCount = 0
	rLock sync.Mutex
)

func renewRequestCount() {
	rLock.Lock()
	requestCount++
	rLock.Unlock()
}

type DevicePluginServer struct {
  deviceinfo.UnimplementedDeviceInfoServer
}

func (DevicePluginServer) ListAllMigDeviceInfos(ctx context.Context, emp *empty.Empty) (*deviceinfo.PodDeviceInfos, error) {
	podDeviceInfos := getMigDeviceInfos()
	return &podDeviceInfos, nil
}

func (DevicePluginServer) WatchMigDeviceInfos(emp *empty.Empty, stream grpc.ServerStreamingServer[deviceinfo.PodDeviceInfos]) (error) {
	for {
		podDeviceInfos := getMigDeviceInfos()
		if err := stream.Send(&podDeviceInfos); err != nil {

			return err
		}
		time.Sleep(5 * time.Second)
	}
	return nil
}

func getMigDeviceInfos()deviceinfo.PodDeviceInfos {
	renewRequestCount()
   	podDeviceInfos := deviceinfo.PodDeviceInfos{
   		PodDeviceinfos: map[string]*deviceinfo.PodDeviceInfo{},
	}
   	podAmount := 3
   	i := 0
	for {
		if i >= podAmount {
			break
		}
		podDeviceInfo := deviceinfo.PodDeviceInfo{
			PodName: "Pod" + strconv.Itoa(i),
		}
		podDeviceInfo.ContainerDevices = map[string]*deviceinfo.ContainerDevices{}
		j := 0
		containerAmount := 2
		for {

			if j >= containerAmount {
				break
			}
			containerDevices := &deviceinfo.ContainerDevices{
				ContainerName: "Container" + strconv.Itoa(j),
                ContainerDevices: []*deviceinfo.ContainerDevice{},
			}

			k := 0
			cardAmount := 2
			for {
                if cardAmount <=  k {
                	break
				}
				containerDevices.ContainerDevices = append(containerDevices.ContainerDevices, &deviceinfo.ContainerDevice{
					ParetDeviceId: "GPU-" + strconv.Itoa(k),
					MigDeviceId: "Mig-" + strconv.Itoa(k) + "-renew-" + strconv.Itoa(requestCount),
				})
				k++
			}

			podDeviceInfo.ContainerDevices[containerDevices.ContainerName] = containerDevices
			j ++

		}
		podDeviceInfos.PodDeviceinfos[podDeviceInfo.PodName] = &podDeviceInfo
		i++
		break
	}
    return podDeviceInfos
}