package postgres

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spaolacci/murmur3"
)

// ShardFn определяет функцию для вычисления номера бакета по ключу.
type ShardFn func(int64) int64

// GetMurmur3Hashing возвращает функцию хеширования Murmur3 для распределения по бакетам.
func GetMurmur3Hashing(buckets int64) ShardFn {
	return func(key int64) int64 {
		hash := int64(murmur3.Sum64([]byte(fmt.Sprintf("%d", key)))) // #nosec G115
		idx := hash % buckets
		if idx < 0 {
			idx += buckets
		}
		return idx
	}
}

// ShardManager управляет шардами и их пулами соединений.
type ShardManager struct {
	hashFn  ShardFn
	buckets int64
	shards  []*Shard

	pools []*pgxpool.Pool
}

// Shard представляет отдельный шард с пулом соединений и позицией последнего бакета (excluded) на окружности хэширования.
type Shard struct {
	Pool           *pgxpool.Pool
	BucketPosition int64 // excluded
}

// NewShardManager создает новый менеджер шардов.
func NewShardManager(hashFn ShardFn, buckets int64, shards []*Shard) (*ShardManager, error) {
	for i := 1; i < len(shards); i++ {
		if shards[i].BucketPosition <= shards[i-1].BucketPosition {
			return nil, errors.New("шарды должны быть отсортированы в порядке возрастания по bucketPosition")
		}
	}

	pools := make([]*pgxpool.Pool, 0, len(shards))

	for _, shard := range shards {
		if shard.BucketPosition > buckets {
			return nil, errors.New("bucketPosition каждого шарда должен быть меньше или равен buckets")
		}
		pools = append(pools, shard.Pool)
	}

	return &ShardManager{
		hashFn:  hashFn,
		buckets: buckets,
		shards:  shards,
		pools:   pools,
	}, nil
}

// GetShardPool возвращает пул соединений для ключа и номер бакета.
func (s *ShardManager) GetShardPool(key int64) (*pgxpool.Pool, int64) {
	bucketIdx := s.hashFn(key)

	return s.getPoolByBucketID(bucketIdx), bucketIdx
}

// GetShardPoolByID возвращает пул соединений по идентификатору бакета в id сущности.
func (s *ShardManager) GetShardPoolByID(id int64) *pgxpool.Pool {
	bucketIdx := id % s.buckets

	return s.getPoolByBucketID(bucketIdx)
}

func (s *ShardManager) getPoolByBucketID(bucketIdx int64) *pgxpool.Pool {
	for _, shard := range s.shards {
		if bucketIdx < shard.BucketPosition {
			return shard.Pool
		}
	}

	return s.shards[0].Pool
}

// GetAllPools возвращает все пулы соединений.
func (s *ShardManager) GetAllPools() []*pgxpool.Pool {
	return s.pools
}
