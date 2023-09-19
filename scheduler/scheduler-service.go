package scheduler

import (
	"fmt"

	"github.com/go-co-op/gocron"
	logUtils "github.com/itbasis/go-log-utils/v2"
	"github.com/juju/zaputil/zapctx"
	"go.uber.org/zap"
)

type FuncJob = func()
type FuncCustomizeSchedule = func(scheduler *gocron.Scheduler) *gocron.Scheduler

type Service interface {
	Schedule(*gocron.Scheduler)
	Job()
}

type AbstractScheduler struct {
	Service

	logger *zap.SugaredLogger

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
	logger := zapctx.Default.With(zap.String(logUtils.MdcSchedulerName, schedulerName)).Sugar()

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

func (receiver *AbstractScheduler) GetLogger() *zap.SugaredLogger {
	return receiver.logger
}

// Schedule
//
// issue: Add prometheus metrics https://github.com/go-co-op/gocron/issues/317
func (receiver *AbstractScheduler) Schedule(scheduler *gocron.Scheduler) {
	if _, err := receiver.funcCustomizeSchedule(scheduler).Do(receiver.funcJob); err != nil {
		receiver.logger.Error(fmt.Errorf("failed to start job: %w", err))
	}
}

func (receiver *AbstractScheduler) defaultCustomizeSchedule(scheduler *gocron.Scheduler) *gocron.Scheduler {
	receiver.logger.Warn("The default settings are used for the scheduler")

	return scheduler.Every(DefaultEveryInterval).Seconds().WaitForSchedule()
}
