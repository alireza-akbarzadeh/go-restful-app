package repository

import (
	"context"

	appErrors "github.com/alireza-akbarzadeh/ginflow/internal/errors"
	"gorm.io/gorm"
)

// TxManager manages database transactions
type TxManager struct {
	db *gorm.DB
}

// NewTxManager creates a new transaction manager
func NewTxManager(db *gorm.DB) *TxManager {
	return &TxManager{db: db}
}

// WithTx executes a function within a database transaction
func (tm *TxManager) WithTx(ctx context.Context, fn func(ctx context.Context, tx *gorm.DB) error) error {
	tx := tm.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return appErrors.New(appErrors.ErrDatabaseOperation, "failed to begin transaction")
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(ctx, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return appErrors.New(appErrors.ErrDatabaseOperation, "failed to commit transaction")
	}

	return nil
}

// TxKey is the key for storing transaction in context
type TxKey struct{}

// WithTxContext adds transaction to context
func WithTxContext(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, TxKey{}, tx)
}

// TxFromContext gets transaction from context
func TxFromContext(ctx context.Context) (*gorm.DB, bool) {
	tx, ok := ctx.Value(TxKey{}).(*gorm.DB)
	return tx, ok
}

// GetDB returns the transaction from context or the default DB
func (tm *TxManager) GetDB(ctx context.Context) *gorm.DB {
	if tx, ok := TxFromContext(ctx); ok {
		return tx
	}
	return tm.db.WithContext(ctx)
}
