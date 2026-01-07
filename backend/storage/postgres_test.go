package storage

import (
	"context"
	"testing"
	"time"

	"github.com/AndreiAlbert/sysmnt/pb"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestPostgresStore_Integration(t *testing.T) {
	ctx := context.Background()
	pgContainer, err := postgres.Run(ctx, "postgres:15-alpine", postgres.WithDatabase("testdb"), postgres.WithUsername("user"), postgres.WithPassword("password"), testcontainers.WithWaitStrategy(
		wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second),
	))
	if err != nil {
		t.Fatalf("Failed to start container: %s", err)
	}

	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("Failed to terminate container: %v", err)
		}
	}()

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}

	store, err := NewPostgresStore(connStr)
	assert.NoError(t, err, "NewPostgresStore should not return an error")
	assert.NotNil(t, store)

	t.Run("Save and GetHistory", func(t *testing.T) {
		stats := pb.SystemStats{
			InstanceId: "test",
			RamUsage:   4096,
			CpuUsage:   25.3,
		}

		err := store.Save(ctx, &stats)
		assert.NoError(t, err, "Save should not return an error")

		history, err := store.GetHistory(ctx)
		assert.NoError(t, err, "GetHistory should not return an errorr")
		assert.NotEmpty(t, history, "History should not be empty")

		found := false
		for _, item := range history {
			if item.InstanceId == "test" {
				assert.Equal(t, item.CpuUsage, 25.3)
				assert.Equal(t, item.RamUsage, float64(4096))
				found = true
			}
		}
		assert.True(t, found, "The saved record was not found in history")
	})
}
