package postgres

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spaolacci/murmur3"
)

type ShardFn func(int64) int64

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

type ShardManager struct {
	hashFn  ShardFn
	buckets int64
	shards  []*Shard

	pools []*pgxpool.Pool
}

type Shard struct {
	Pool           *pgxpool.Pool
	BucketPosition int64 // excluded
}

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

func (s *ShardManager) GetShardPool(key int64) (*pgxpool.Pool, int64) {
	bucketIdx := s.hashFn(key)

	return s.getPoolByBucketID(bucketIdx), bucketIdx
}

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

func (s *ShardManager) GetAllPools() []*pgxpool.Pool {
	return s.pools
}
