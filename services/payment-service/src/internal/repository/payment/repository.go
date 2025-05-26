package payment

import (
	gorm2 "gorm.io/gorm"
	"payment-service/internal/model"
	"pkg/infrastructure/DB/gorm"
)

type Repository struct {
	db *gorm.DBService
}

func (r Repository) query() (tx *gorm2.DB) {
	return r.db.ActiveTx().Model(&model.Payment{})
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
