package metrics

import (
	"route256/cart/pkg/metrics"
	"time"
)

// RepositoryInfo содержит информацию для опроса репозитория
type RepositoryInfo struct {
	Repo        repository
	ObjectsName string
	Interval    time.Duration
}

// RepositoryObserver собирает метрики с репозиториев
type RepositoryObserver struct {
	repos []*RepositoryInfo
	done  chan struct{}
}

type repository interface {
	CountObjects() int
}

// NewRepositoryObserver создает новый RepositoryObserver и запускает сбор метрик с репозиториев
func NewRepositoryObserver(repos []*RepositoryInfo) *RepositoryObserver {
	r := &RepositoryObserver{
		repos: repos,
	}

	for _, repoInfo := range repos {
		go func() {
			t := time.NewTicker(repoInfo.Interval)
			defer t.Stop()

			for {
				select {
				case <-t.C:
					metrics.StoreRepositorySize(repoInfo.ObjectsName, float64(repoInfo.Repo.CountObjects()))
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
