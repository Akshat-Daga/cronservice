// pkg/crontab/scheduler.go
package crontab

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
)

// Scheduler is a thin wrapper around go-co-op/gocron scheduler.
type Scheduler struct {
	s       *gocron.Scheduler
	elector gocron.Elector
	locker  gocron.Locker
}

// NewScheduler creates a local scheduler (non-distributed).
func NewScheduler() *Scheduler {
	s := gocron.NewScheduler(time.Local)
	return &Scheduler{s: s}
}

// NewSchedulerWithDistributed creates a scheduler with optional elector/locker
// (note: how you wire them depends on your gocron version and infra).
func NewSchedulerWithDistributed(elector gocron.Elector, locker gocron.Locker) *Scheduler {
	s := gocron.NewScheduler(time.Local)
	// NOTE: gocron/v2 supports SetElector/SetLocker on the scheduler instance
	// in some versions. If your version supports it, you can wire them here.
	// Example (uncomment if supported):
	// if elector != nil {
	//     s.SetElector(elector)
	// }
	// if locker != nil {
	//     s.SetLocker(locker)
	// }
	return &Scheduler{s: s, elector: elector, locker: locker}
}

// EveryInterval schedules a function every <interval> <unit>.
// unit: second(s), minute(s), hour(s), day(s), week(s)
func (sch *Scheduler) EveryInterval(interval uint64, unit string, jobFunc func()) error {
	if interval == 0 {
		return fmt.Errorf("interval must be > 0")
	}

	switch unit {
	case "second", "seconds":
		_, err := sch.s.Every(interval).Second().Do(jobFunc)
		return err
	case "minute", "minutes":
		_, err := sch.s.Every(interval).Minute().Do(jobFunc)
		return err
	case "hour", "hours":
		_, err := sch.s.Every(interval).Hour().Do(jobFunc)
		return err
	case "day", "days":
		_, err := sch.s.Every(interval).Day().Do(jobFunc)
		return err
	case "week", "weeks":
		_, err := sch.s.Every(interval).Week().Do(jobFunc)
		return err
	default:
		return fmt.Errorf("invalid unit: %s", unit)
	}
}

// Cron schedules a cron expression (standard cron-like) running the jobFunc.
func (sch *Scheduler) Cron(cronExpr string, jobFunc func()) error {
	_, err := sch.s.Cron(cronExpr).Do(jobFunc)
	if err != nil {
		return fmt.Errorf("failed to schedule cron job: %v", err)
	}
	return nil
}

// StartAsync runs the scheduler in a goroutine.
func (sch *Scheduler) StartAsync() {
	sch.s.StartAsync()
}

// Stop stops the scheduler (graceful shutdown).
func (sch *Scheduler) Stop() {
	sch.s.Stop()
}

// Jobs returns a slice of scheduled jobs (read-only).
func (sch *Scheduler) Jobs() []*gocron.Job {
	return sch.s.Jobs()
}

// Elector and Locker accessors (may be nil if not set).
func (sch *Scheduler) Elector() gocron.Elector {
	return sch.elector
}

func (sch *Scheduler) Locker() gocron.Locker {
	return sch.locker
}
