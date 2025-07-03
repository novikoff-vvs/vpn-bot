package plan

import (
	"user-service/internal/models"
	"user-service/internal/repository/plan/sqlite"
)

var upgradeMap = map[string]string{
	models.SubscriptionTestCode: models.SubscriptionBaseCode,
}

type Service struct {
	repo *sqlite.PlanRepository
}

func NewPlanService(repo *sqlite.PlanRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetUpgradedPlan(plan models.Plan) (models.Plan, error) {
	var newCode string
	var ok bool

	if newCode, ok = upgradeMap[plan.Code]; !ok {
		return plan, nil
	}
	plan, err := s.repo.GetByCode(newCode)
	if err != nil {
		return models.Plan{}, err
	}
	return plan, nil
}
