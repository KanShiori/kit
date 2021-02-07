package timer

import (
	"fmt"
	"time"
)

type Timer interface {
	// Clock 检查当前时间是否允许, 允许情况下会执行需要的置位(TimeSpan)
	Clock() (bool, time.Time)

	// Clock 检查指定时间是否允许, 允许情况下会执行需要的置位(TimeSpan)
	ClockIn(tm time.Time) (bool, time.Time)

	// IsTimeUpNow 检查指定时间是否允许
	IsTimeUpNow() (bool, time.Time)

	// IsTimeUp 检查指定时间是否允许
	IsTimeUp(tm time.Time) (bool, time.Time)
}

// TimerList 用于检测多个 Timer 是否同时 time up
//
// Note: next 会有误差, 因为只是取的 timer 中 next 最大值
type TimerList struct {
	timers []Timer
}

func NewTimerList(timers ...Timer) (*TimerList, error) {
	if len(timers) == 0 {
		return nil, fmt.Errorf("empty timer")
	}

	return &TimerList{
		timers: timers,
	}, nil
}

func (t *TimerList) Clock() (bool, time.Time) {
	return t.ClockIn(time.Now())
}

func (t *TimerList) ClockIn(tm time.Time) (bool, time.Time) {
	is, next := t.IsTimeUp(tm)
	if !is {
		return is, next
	}

	for _, timer := range t.timers {
		timer.ClockIn(tm)
	}

	return is, next
}

func (t *TimerList) IsTimeUpNow() (bool, time.Time) {
	return t.IsTimeUp(time.Now())
}

func (t *TimerList) IsTimeUp(tm time.Time) (bool, time.Time) {
	maxnext := time.Time{}
	iss := true

	for _, timer := range t.timers {
		is, next := timer.IsTimeUp(tm)
		if !is {
			iss = false
			if next.After(maxnext) {
				maxnext = next
			}
		}
	}

	if !iss {
		return false, maxnext
	}
	return true, tm
}
