package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/AndreiAlbert/sysmnt/pb"
	"github.com/AndreiAlbert/sysmnt/storage"
)

type server struct {
	pb.UnimplementedMonitorServiceServer
	hub   *StreamHub
	store storage.Storage
}

func (s *server) PushStats(stream pb.MonitorService_PushStatsServer) error {
	for {
		stats, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Empty{})
		}
		if err != nil {
			log.Printf("gRPC Stream Error: %v", err)
			return err
		}
		stats.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
		go func(data *pb.SystemStats) {
			if err := s.store.Save(context.Background(), data); err != nil {
				log.Printf("Failed to save data: %v", err)
			}
		}(stats)
		log.Print("the agent sent data")
		s.hub.broadcast <- stats
	}
}
