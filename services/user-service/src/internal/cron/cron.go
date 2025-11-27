package cron

import (
	"user-service/internal/job/subscription"

	"github.com/novikoff-vvs/logger"
	"github.com/robfig/cron/v3"
)

type Worker struct {
	cron *cron.Cron
	lg   logger.Interface
}

func (w Worker) Start() {
	w.cron.Start()
}

func NewWorker(subscriptionExpiringJob *subscription.ExpiringNotificationJob, lg logger.Interface) (*Worker, error) {
	c, err := initCron(subscriptionExpiringJob, lg)
	if err != nil {
		return nil, err
	}

	return &Worker{
		cron: c,
		lg:   lg,
	}, nil
}

func initCron(subscriptionExpiringJob *subscription.ExpiringNotificationJob, lg logger.Interface) (*cron.Cron, error) {
	c := cron.New()

	// Запускаем каждый день в 2:00 ночи
	_, err := c.AddJob("0 2 * * *", subscriptionExpiringJob)
	if err != nil {
		return nil, err
	}

	lg.Info("Cron jobs initialized - subscription expiring notifications will run daily at 2:00 AM")

	return c, nil
}
