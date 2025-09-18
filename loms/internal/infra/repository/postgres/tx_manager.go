package postgres

import (
	"context"
	"fmt"
	"route256/cart/pkg/logger"
	"route256/loms/internal/service"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txKey struct{}

func ctxWithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// TxFromCtx извлекает транзакцию из контекста.
func TxFromCtx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}

// TxManager реализует интерфейс service.TxManager для работы с транзакциями.
type TxManager struct {
	poolManager PoolManager
}

// NewPgTxManager создает новый TxManager.
func NewPgTxManager(poolManager PoolManager) *TxManager {
	return &TxManager{
		poolManager: poolManager,
	}
}

func (m *TxManager) getPool(operationType service.OperationType) *pgxpool.Pool {
	if operationType == service.Read {
		return m.poolManager.Readable()
	}

	return m.poolManager.Writable()
}

// WithTransaction выполняет fn в транзакции с дефолтным уровнем изоляции.
func (m *TxManager) WithTransaction(ctx context.Context, operationType service.OperationType, fn func(ctx context.Context) error) (err error) {
	return m.WithTx(ctx, operationType, pgx.TxOptions{}, fn)
}

// WithRepeatableRead выполняет fn в транзакции с уровнем изоляции RepeatableRead.
func (m *TxManager) WithRepeatableRead(ctx context.Context, operationType service.OperationType, fn func(ctx context.Context) error) (err error) {
	return m.WithTx(ctx, operationType, pgx.TxOptions{IsoLevel: pgx.RepeatableRead}, fn)
}

// WithTx выполняет fn в транзакции.
func (m *TxManager) WithTx(ctx context.Context, operationType service.OperationType, options pgx.TxOptions, fn func(ctx context.Context) error) (err error) {
	var tx pgx.Tx

	existedTx, existTx := TxFromCtx(ctx)
	if existTx {
		tx, err = existedTx.Begin(ctx)
		if err != nil {
			return
		}
	} else {
		tx, err = m.getPool(operationType).BeginTx(ctx, options)
		if err != nil {
			return
		}
		ctx = ctxWithTx(ctx, tx)
	}

	defer func() {
		if p := recover(); p != nil {
			roolbackErr := tx.Rollback(ctx)
			if roolbackErr != nil {
				logger.Warning(fmt.Sprintf("tx.Rollback: %s", roolbackErr.Error()))
			}
			panic(p)
		} else if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				logger.Warning(fmt.Sprintf("tx.Rollback: %s", rollbackErr.Error()))
			}
		} else {
			err = tx.Commit(ctx)
		}
	}()

	err = fn(ctx)

	return
}
