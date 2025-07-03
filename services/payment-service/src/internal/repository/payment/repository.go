package payment

import (
	"payment-service/internal/model"
	"pkg/infrastructure/DB/gorm"
)

type Repository struct {
	db *gorm.DBService
}

func (r Repository) LogPayment(m model.Payment) error {
	if err := r.db.Begin().Create(&m).Error; err != nil {
		err = r.db.Rollback()
		if err != nil {
			return err
		}
		return err
	}
	err := r.db.Commit()
	if err != nil {
		return err
	}
	return nil
}

func NewPaymentRepository(db *gorm.DBService) *Repository {
	return &Repository{db: db}
}
