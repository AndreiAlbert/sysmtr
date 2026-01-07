package main

import (
	"context"
	"io"
	"sync"
	"testing"

	"github.com/AndreiAlbert/sysmnt/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Save(ctx context.Context, stats *pb.SystemStats) error {
	args := m.Called(ctx, stats)
	return args.Error(0)
}

func (m *MockStorage) GetHistory(ctx context.Context) ([]*pb.SystemStats, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*pb.SystemStats), args.Error(1)
}

type MockStream struct {
	grpc.ServerStream
	Messages []*pb.SystemStats
	Index    int
}

func (m *MockStream) Recv() (*pb.SystemStats, error) {
	if m.Index >= len(m.Messages) {
		return nil, io.EOF
	}
	msg := m.Messages[m.Index]
	m.Index++
	return msg, nil
}

func (m *MockStream) SendAndClose(e *pb.Empty) error {
	return nil
}

func TestPushStats(t *testing.T) {
	var wg sync.WaitGroup
	expectedCalls := 2
	wg.Add(expectedCalls)
	mockStore := new(MockStorage)
	mockStore.On("Save", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			wg.Done()
		}).Return(nil)
	hub := NewStreamHub()
	go hub.run()
	s := &server{
		hub:   hub,
		store: mockStore,
	}
	mockStream := MockStream{
		Messages: []*pb.SystemStats{
			{InstanceId: "test-1", CpuUsage: 10.0},
			{InstanceId: "test-2", CpuUsage: 20.0},
		},
	}
	err := s.PushStats(&mockStream)
	assert.NoError(t, err)
	wg.Wait()
	mockStore.AssertNumberOfCalls(t, "Save", 2)
}
