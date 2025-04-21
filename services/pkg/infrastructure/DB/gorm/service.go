package gorm

import (
	"gorm.io/gorm"
	"sync"
)

type DBService struct {
	db       *gorm.DB
	activeTx *TxWrapper
	mu       sync.Mutex
}

func NewDBService(db *gorm.DB) *DBService {
	return &DBService{db: db}
}

// Начинает транзакцию, если нет активной
func (s *DBService) Begin() *gorm.DB {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.activeTx != nil && s.activeTx.IsActive() {
		return s.activeTx.GetDB()
	}

	tx := s.db.Begin()
	s.activeTx = NewTxWrapper(tx)
	return tx
}

// Возвращает активную транзакцию или nil
func (s *DBService) ActiveTx() *gorm.DB {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.activeTx != nil && s.activeTx.IsActive() {
		return s.activeTx.GetDB()
	}
	return nil
}

// Коммитит активную транзакцию (если есть) и сбрасывает
func (s *DBService) Commit() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.activeTx == nil || !s.activeTx.IsActive() {
		return nil
	}

	err := s.activeTx.Commit()
	s.activeTx = nil
	return err
}

// Откатывает активную транзакцию (если есть) и сбрасывает
func (s *DBService) Rollback() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.activeTx == nil || !s.activeTx.IsActive() {
		return nil
	}

	err := s.activeTx.Rollback()
	s.activeTx = nil
	return err
}

// Доступ к сырой DB (без транзакции)
func (s *DBService) DB() *gorm.DB {
	return s.db
}
