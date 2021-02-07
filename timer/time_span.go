package timer

import (
	"fmt"
	"time"
)

type TimeSpan struct {
	interval time.Duration
	lastAt   time.Time
}

func NewTimeSpan(interval time.Duration, lastAt time.Time) *TimeSpan {
	return &TimeSpan{
		interval: interval,
		lastAt:   time.Time{},
	}
}

// Clock 与 now 比较, 当timeup时置 last 为  now
func (t *TimeSpan) Clock() (bool, time.Time) {
	return t.ClockIn(time.Now())
}

// ClockIn 与 ts 比较, 当 timeup 时置 last 为 tm
func (t *TimeSpan) ClockIn(tm time.Time) (bool, time.Time) {
	is, next := t.IsTimeUp(tm)
	if is {
		t.ResetAs(tm)
	}

	return is, next
}

// IsTimeUpNow 与 now_ts 比较
func (t *TimeSpan) IsTimeUpNow() (bool, time.Time) {
	return t.IsTimeUp(time.Now())
}

// IsTimeUp 与 ts 比较
func (t *TimeSpan) IsTimeUp(tm time.Time) (bool, time.Time) {
	next := t.lastAt.Add(t.interval)
	fmt.Println(next, tm)
	if tm.Before(next) {
		return false, next
	}
	return true, tm
}

// ResetNow 置 last 为  now_ts
func (t *TimeSpan) ResetNow() {
	t.ResetAs(time.Now())
}

// ResetNow 置 last 为  ts
func (t *TimeSpan) ResetAs(tm time.Time) {
	t.lastAt = tm
}

func (t *TimeSpan) LastAt() time.Time {
	return t.lastAt
}

func (t *TimeSpan) Interval() time.Duration {
	return t.interval
}
