package sqlite

import (
	"pkg/infrastructure/DB/gorm"
	"user-service/internal/models"
)

type PlanRepository struct {
	dbService *gorm.DBService
}

func (r *PlanRepository) BeginTransaction() {
	r.dbService.Begin()
}

func (r *PlanRepository) CommitTransaction() error {
	return r.dbService.Commit()
}

func (r *PlanRepository) RollbackTransaction() error {
	return r.dbService.Rollback()
}

func (r *PlanRepository) GetByCode(code string) (models.Plan, error) {
	var result models.Plan
	err := r.dbService.DB().Model(&result).Where("code=?", code).First(&result).Error
	if err != nil {
		return models.Plan{}, err
	}
	return result, nil
}

func NewPlanRepository(dbService *gorm.DBService) *PlanRepository {
	return &PlanRepository{dbService: dbService}
}
