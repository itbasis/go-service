package scheduler

import (
	"github.com/go-co-op/gocron"
	logUtils "github.com/itbasis/go-log-utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type FuncJob = func()
type FuncCustomizeSchedule = func(scheduler *gocron.Scheduler) *gocron.Scheduler

type Service interface {
	Schedule(*gocron.Scheduler)
	Job()
}

type AbstractScheduler struct {
	Service

	logger zerolog.Logger

	funcJob               func()
	funcCustomizeSchedule FuncCustomizeSchedule
}

func NewAbstractScheduler(schedulerName string, funcJob func()) *AbstractScheduler {
	return NewAbstractSchedulerWithCustomizeSchedule(schedulerName, funcJob, nil)
}

func NewAbstractSchedulerWithCustomizeSchedule(
	schedulerName string,
	funcJob func(),
	funcCustomizeSchedule *FuncCustomizeSchedule,
) *AbstractScheduler {
	logger := log.Logger.With().Str(logUtils.MdcSchedulerName, schedulerName).Logger()
	abstractScheduler := AbstractScheduler{
		logger:  logger,
		funcJob: funcJob,
	}

	if funcCustomizeSchedule == nil {
		abstractScheduler.funcCustomizeSchedule = abstractScheduler.defaultCustomizeSchedule
	} else {
		abstractScheduler.funcCustomizeSchedule = *funcCustomizeSchedule
	}

	return &abstractScheduler
}

func (receiver *AbstractScheduler) GetLogger() zerolog.Logger {
	return receiver.logger
}

func (receiver *AbstractScheduler) Schedule(scheduler *gocron.Scheduler) {
	// issue: Add prometheus metrics https://github.com/go-co-op/gocron/issues/317

	if _, err := receiver.funcCustomizeSchedule(scheduler).Do(receiver.funcJob); err != nil {
		receiver.logger.Error().Err(err).Msg("Failed to start job")
	}
}

func (receiver *AbstractScheduler) defaultCustomizeSchedule(scheduler *gocron.Scheduler) *gocron.Scheduler {
	receiver.logger.Warn().Msg("The default settings are used for the scheduler")

	return scheduler.Every(5).Seconds().WaitForSchedule()
}
