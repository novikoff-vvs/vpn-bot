package gorm

import (
	"gorm.io/gorm"
	"pkg/exceptions"
)

type TxWrapper struct {
	tx         *gorm.DB
	committed  bool
	rolledBack bool
}

func NewTxWrapper(tx *gorm.DB) *TxWrapper {
	return &TxWrapper{tx: tx}
}

func (tw *TxWrapper) Commit() error {
	err := tw.status()
	if err != nil {
		return err
	}
	tw.committed = true
	return tw.tx.Commit().Error
}

func (tw *TxWrapper) Rollback() error {
	err := tw.status()
	if err != nil {
		return err
	}
	tw.rolledBack = true
	return tw.tx.Rollback().Error
}

func (tw *TxWrapper) GetDB() *gorm.DB {
	return tw.tx
}

func (tw *TxWrapper) IsActive() bool {
	return !tw.committed && !tw.rolledBack
}

func (tw *TxWrapper) status() error {
	if tw.committed {
		return exceptions.ErrCommitedAlready
	}
	if tw.rolledBack {
		return exceptions.ErrRollbackAlready
	}
	return nil
}
