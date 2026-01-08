package main

import (
	"context"
	"log"
	"time"

	"github.com/AndreiAlbert/sysmnt/config"
	"github.com/AndreiAlbert/sysmnt/pb"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.LoadConfig()
	conn, err := grpc.Dial(cfg.ServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewMonitorServiceClient(conn)
	stream, err := client.PushStats(context.Background())
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		v, _ := mem.VirtualMemory()
		c, _ := cpu.Percent(0, false)
		stats := &pb.SystemStats{
			CpuUsage:   c[0],
			RamUsage:   v.UsedPercent,
			InstanceId: "localhost",
		}
		if err := stream.Send(stats); err != nil {
			log.Printf("Failed to send stats: %v", err)
			break
		}
	}

}
