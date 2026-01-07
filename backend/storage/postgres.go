package storage

import (
	"context"
	"database/sql"

	"github.com/AndreiAlbert/sysmnt/pb"
	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(connStr string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	query := ` 
	CREATE TABLE IF NOT EXISTS system_stats(
		id SERIAL PRIMARY KEY,
		instance_id TEXT NOT NULL, 
		cpu_usage FLOAT NOT NULL,
		ram_usage FLOAT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(query); err != nil {
		return nil, err
	}
	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Save(ctx context.Context, stats *pb.SystemStats) error {
	query := `INSERT INTO system_stats (instance_id, cpu_usage, ram_usage) VALUES ($1, $2, $3)`
	_, err := s.db.ExecContext(ctx, query, stats.InstanceId, stats.CpuUsage, stats.RamUsage)
	return err
}

func (s *PostgresStore) GetHistory(ctx context.Context) ([]*pb.SystemStats, error) {
	query := `SELECT instance_id, cpu_usage, ram_usage, TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI:SS') as created_at FROM system_stats ORDER BY created_at DESC LIMIT 50`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*pb.SystemStats
	for rows.Next() {
		var stats pb.SystemStats
		if err := rows.Scan(&stats.InstanceId, &stats.CpuUsage, &stats.RamUsage, &stats.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, &stats)
	}
	return result, nil
}
