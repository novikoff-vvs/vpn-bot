package payment

import (
	"github.com/google/uuid"
	"payment-service/internal/model"
	"payment-service/internal/repository/payment"
)

type Service struct {
	repo *payment.Repository
}

func (s Service) LogPayment(m model.Payment) error {
	m.UUID = uuid.New().String()
	err := s.repo.LogPayment(m)
	if err != nil {
		return err
	}
	return nil
}

func NewPaymentService(repo *payment.Repository) *Service {
	return &Service{
		repo: repo,
	}
}
