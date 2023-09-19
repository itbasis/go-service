package service

import (
	"github.com/go-co-op/gocron"
	"github.com/itbasis/go-service/v2/utils/time"
	"github.com/juju/zaputil/zapctx"
)

func (receiver *Service) GetScheduler() *gocron.Scheduler {
	if receiver.scheduler == nil && receiver.config.SchedulerEnabled {
		receiver.scheduler = gocron.NewScheduler(time.GlobalTime)
	}

	if receiver.scheduler == nil {
		zapctx.Default.Sugar().Warn("scheduler is not enabled")
	}

	return receiver.scheduler
}
