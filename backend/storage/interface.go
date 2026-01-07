package storage

import (
	"context"

	"github.com/AndreiAlbert/sysmnt/pb"
)

type Storage interface {
	Save(ctx context.Context, stats *pb.SystemStats) error
	GetHistory(ctx context.Context) ([]*pb.SystemStats, error)
}
