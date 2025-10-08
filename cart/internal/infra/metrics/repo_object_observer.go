package metrics

import (
	"route256/cart/pkg/metrics"
	"time"
)

// RepositoryInfo содержит информацию для опроса репозитория
type RepositoryInfo struct {
	Repo       repository
	ObjectName string
}

// RepositoryObserver собирает метрики с репозиториев
type RepositoryObserver struct {
	repos    []*RepositoryInfo
	done     chan struct{}
	interval time.Duration
}

type repository interface {
	CountObjects() int
}

// NewRepositoryObserver создает новый RepositoryObserver и запускает сбор метрик с репозиториев
func NewRepositoryObserver(repos []*RepositoryInfo, interval time.Duration) *RepositoryObserver {
	r := &RepositoryObserver{
		repos:    repos,
		done:     make(chan struct{}),
		interval: interval,
	}

	for _, repoInfo := range repos {
		go func() {
			t := time.NewTicker(r.interval)
			defer t.Stop()

			for {
				select {
				case <-t.C:
					metrics.StoreRepositorySize(repoInfo.ObjectName, float64(repoInfo.Repo.CountObjects()))
				case <-r.done:
					return
				}
			}
		}()
	}

	return r
}

// Stop останавливает сбор метрик с репозиториев
func (r *RepositoryObserver) Stop() {
	close(r.done)
}
