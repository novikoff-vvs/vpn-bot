package cron

import (
	"github.com/novikoff-vvs/logger"
	"github.com/robfig/cron/v3"
	"pkg/infrastructure/client/user"
	userJob "vpn-service/internal/job/user"
	"vpn-service/internal/service/vpn"
)

type Worker struct {
	cron *cron.Cron
	lg   logger.Interface
}

func (w Worker) Start() {
	w.cron.Start()
}

func NewWorker(serviceInterface vpn.ServiceInterface, client *user.Client, lg logger.Interface) (*Worker, error) {
	c, err := initCron(serviceInterface, client, lg)
	if err != nil {
		return nil, err
	}

	return &Worker{
		cron: c,
		lg:   lg,
	}, nil
}

func initCron(serviceInterface vpn.ServiceInterface, client *user.Client, lg logger.Interface) (*cron.Cron, error) {
	c := cron.New()
	job1 := userJob.NewSyncJob(client, serviceInterface, lg)
	_, err := c.AddJob("* * * * *", job1)
	if err != nil {
		return nil, err
	}

	return c, nil
}
