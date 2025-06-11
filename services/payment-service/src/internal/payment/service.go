package payment

import (
	"github.com/google/uuid"
	"payment-service/internal/model"
	"payment-service/internal/repository/payment"
	"pkg/infrastructure/client/subscription"
)

type Service struct {
	repo         *payment.Repository
	subscrClient *subscription.Client
}

func (s Service) LogPayment(m model.Payment) error {
	m.UUID = uuid.New().String()
	err := s.repo.LogPayment(m)
	if err != nil {
		return err
	}
	_, err = s.subscrClient.RefreshSubscription(subscription.RefreshRequest{
		UserUUID:      m.Token,
		AmountPeriods: 1,
	})
	if err != nil {
		return err
	}
	return nil
}

func NewPaymentService(repo *payment.Repository, subscrClient *subscription.Client) *Service {
	return &Service{
		repo:         repo,
		subscrClient: subscrClient,
	}
}
