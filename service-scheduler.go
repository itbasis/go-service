package service

import (
	"github.com/go-co-op/gocron"
	"github.com/itbasis/go-service/utils/time"
	"github.com/rs/zerolog/log"
)

func (receiver *Service) GetScheduler() *gocron.Scheduler {
	if receiver.scheduler == nil && receiver.config.SchedulerEnabled {
		receiver.scheduler = gocron.NewScheduler(time.GlobalTime)
	}

	if receiver.scheduler == nil {
		log.Warn().Msg("scheduler is not enabled")
	}

	return receiver.scheduler
}
