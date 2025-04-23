package cron

import (
	"github.com/robfig/cron/v3"
	"user-service/job"
)

type Worker struct {
	cron *cron.Cron
}

func (w Worker) Start() {
	w.cron.Start()
}

func NewWorker() (*Worker, error) {
	c, err := initCron()
	if err != nil {
		return nil, err
	}

	return &Worker{
		cron: c,
	}, nil
}

func initCron() (*cron.Cron, error) {
	c := cron.New()
	_, err := c.AddJob("@hourly", job.SyncUsersJob{})
	if err != nil {
		return nil, err
	}

	return c, nil
}
