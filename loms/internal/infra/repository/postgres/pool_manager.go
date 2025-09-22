package postgres

import (
	"sync/atomic"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RRPoolManager управляет пулом соединений с master и replica базами данных.
// Распределяет нагрузку по нодам алгоритмом Round-Robin.
type RRPoolManager struct {
	writable    *pgxpool.Pool   // пул для записи (master)
	readable    []*pgxpool.Pool // пулы для чтения (master + replicas)
	readableIdx uint64          // индекс для round-robin
}

// NewRRPoolManager создает новый RRPoolManager с master и replica пулами.
func NewRRPoolManager(master *pgxpool.Pool, replicas []*pgxpool.Pool) (*RRPoolManager, error) {
	poolManager := &RRPoolManager{
		writable: master,
		readable: make([]*pgxpool.Pool, 0, len(replicas)+1),
	}
	poolManager.readable = append(poolManager.readable, master)
	poolManager.readable = append(poolManager.readable, replicas...)

	return poolManager, nil
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
