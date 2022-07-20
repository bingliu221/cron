package cron

import (
	"context"
	"errors"
	"time"
)

type TaskFunc func(at time.Time)

type Cron struct {
	task TaskFunc

	seconds  *Round
	minutes  *Round
	hours    *Round
	days     *Round
	months   *Round
	weekdays *Round

	year int

	location *time.Location
}

func buildRound(min, max int, selectFunc IntSelector) *Round {
	selected := make([]int, 0)
	for i := min; i <= max; i++ {
		if selectFunc(i) {
			selected = append(selected, i)
		}
	}
	return NewRound(selected, min, max)
}

func New(task TaskFunc, selector TimeSelector) (*Cron, error) {
	r := &Cron{
		task:     task,
		location: time.Local,
	}

	r.seconds = buildRound(0, 59, selector.Second)
	r.minutes = buildRound(0, 59, selector.Minute)
	r.hours = buildRound(0, 23, selector.Hour)
	r.days = buildRound(1, 31, selector.Day)
	r.months = buildRound(1, 12, selector.Month)
	r.weekdays = buildRound(0, 6, selector.Weekday)

	if selector.Location != nil {
		r.location = selector.Location
	}

	return r, nil
}

func (r *Cron) init() {
	t := time.Now().In(r.location)
	r.seconds.ShiftTo(t.Second())
	r.seconds.TickBack()
	r.minutes.ShiftTo(t.Minute())
	r.hours.ShiftTo(t.Hour())
	r.days.ShiftTo(t.Day())
	r.months.ShiftTo(int(t.Month()))
	// weekdays is only use for filtering
	//r.weekdays.ShiftTo(int(t.Weekday()))
	r.year = t.Year()
}

func (r *Cron) current() time.Time {
	return time.Date(r.year, time.Month(r.months.Value()), r.days.Value(), r.hours.Value(), r.minutes.Value(), r.seconds.Value(), 0, r.location)
}

func (r *Cron) nextTick() time.Time {
	carry := r.seconds.Tick()
	if carry {
		carry = r.minutes.Tick()
	}
	if carry {
		carry = r.hours.Tick()
	}
	if carry {
		for {
			carry = r.days.Tick()
			if carry {
				carry = r.months.Tick()
			}
			if carry {
				r.year += 1
			}
			wd := r.current().Weekday()
			if r.weekdays.Contains(int(wd)) {
				break
			}
		}
	}
	return r.current()
}

func (r *Cron) Run(ctx context.Context) {
	r.init()
	for {
		next := r.nextTick()
		timer, cancel := context.WithDeadline(context.Background(), next)
		select {
		case <-ctx.Done():
			// Cron stopped
			cancel() // cancel next job
			return

		case <-timer.Done():
			err := timer.Err()
			if errors.Is(err, context.DeadlineExceeded) {
				go r.task(next)
			}
		}
	}
}
