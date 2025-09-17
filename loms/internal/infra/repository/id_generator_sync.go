package repository

import (
	"sync/atomic"
)

// IDGeneratorSync структура потокобезопасного генератора последовательных ID
type IDGeneratorSync struct {
	counter int64
}

// NewIDGeneratorSync создает новый генератор последовательных ID
func NewIDGeneratorSync() *IDGeneratorSync {
	return &IDGeneratorSync{}
}

// NextID возвращает уникальный ID
func (g *IDGeneratorSync) NextID() int64 {
	return atomic.AddInt64(&g.counter, 1)
}
