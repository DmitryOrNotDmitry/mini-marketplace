package postgres

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RRPoolManager управляет пулом соединений с master и replica базами данных.
type RRPoolManager struct {
	writable    *pgxpool.Pool   // пул для записи (master)
	readable    []*pgxpool.Pool // пулы для чтения (master + replicas)
	readableIdx uint64          // индекс для round-robin
}

// NewRRPoolManager создает новый RRPoolManager с master и replica пулами.
func NewRRPoolManager(ctx context.Context, masterDSN string, replicasDSNs []string) (*RRPoolManager, error) {
	masterPool, err := newPool(ctx, masterDSN)
	if err != nil {
		return nil, fmt.Errorf("newPool: %w", err)
	}

	poolManager := &RRPoolManager{
		writable: masterPool,
		readable: make([]*pgxpool.Pool, 0, len(replicasDSNs)+1),
	}
	poolManager.readable = append(poolManager.readable, masterPool)

	for _, replicaDSN := range replicasDSNs {
		replicaPool, err := newPool(ctx, replicaDSN)
		if err != nil {
			return nil, fmt.Errorf("newPool: %w", err)
		}

		poolManager.readable = append(poolManager.readable, replicaPool)
	}

	return poolManager, nil
}

func newPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig (dsn=%s): %w", dsn, err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.NewWithConfig (dsn=%s): %w", dsn, err)
	}

	return pool, nil
}

// Readable возвращает пул для чтения (round-robin между master и replicas).
func (pm *RRPoolManager) Readable() *pgxpool.Pool {
	idx := atomic.AddUint64(&pm.readableIdx, 1) % uint64(len(pm.readable))
	return pm.readable[idx]
}

// Writable возвращает пул для записи.
func (pm *RRPoolManager) Writable() *pgxpool.Pool {
	return pm.writable
}
