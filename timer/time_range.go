package timer

import (
	"fmt"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

const (
	// AnyTimeCronSpec 是任何时候都符合的 cronspec
	AnyTimeCronSpec = "* * * * *"
)

var (
	cronParser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
)

// TimeRange 描述一个时间范围, 并判断 time 是否符合 range.
// TimeRange 的判断条件来自于 Cron format, 见 https://en.wikipedia.org/wiki/Cron
//
// Ex:
//  NewHourTimeRange(1,18), 表明 1-18 点都是符合(包括 18:59 也是)
// Note:
//  TimeRange 中都是 [left,right] 的
type TimeRange struct {
	spec     string
	schedule cron.Schedule
}

// NewMinTimeRange is same as 'NewTimeRange("left-right * * * *")'
func NewMinTimeRange(left, right uint) (*TimeRange, error) {
	spec := fmt.Sprintf("%d-%d * * * *", left, right)
	return NewTimeRange(spec)
}

// NewHourTimeRange is same as 'NewTimeRange("* left-right * * *")'
func NewHourTimeRange(left, right uint) (*TimeRange, error) {
	spec := fmt.Sprintf("* %d-%d * * *", left, right)
	return NewTimeRange(spec)
}

// NewDomTimeRange is same as 'NewTimeRange("* * * left-right *")'
func NewDomTimeRange(left, right uint) (*TimeRange, error) {
	spec := fmt.Sprintf("* * %d-%d * *", left, right)
	return NewTimeRange(spec)
}

// NewMonthTimeRange is same as 'NewTimeRange("* * * left-right *")'
func NewMonthTimeRange(left, right uint) (*TimeRange, error) {
	spec := fmt.Sprintf("* * * %d-%d *", left, right)
	return NewTimeRange(spec)
}

// NewDowTimeRange is same as 'NewTimeRange("* * * * left-right")'
func NewDowTimeRange(left, right uint) (*TimeRange, error) {
	spec := fmt.Sprintf("* * * * %d-%d", left, right)
	return NewTimeRange(spec)
}

// NewTimeRange 从 cron spec 创建 TimeRange, 每一必须为 left-right 或 *
//
// ┌───────────── minute (0 - 59)
// │ ┌───────────── hour (0 - 23)
// │ │ ┌───────────── day of the month (1 - 31)
// │ │ │ ┌───────────── month (1 - 12)
// │ │ │ │ ┌───────────── day of the week (0 - 6) (Sunday to Saturday;
// │ │ │ │ │                                   7 is also Sunday on some systems)
// │ │ │ │ │
// │ │ │ │ │
// * * * * *
func NewTimeRange(spec string) (*TimeRange, error) {
	// - validate
	{
		fields := strings.Fields(spec)
		if len(fields) != 5 {
			return nil, fmt.Errorf("invaild spec %s", spec)
		}

		for _, field := range fields {
			if field == "*" {
				continue
			}
			if strings.ContainsRune(field, '-') {
				continue
			}
			return nil, fmt.Errorf("invaild spec %s", spec)
		}

	}

	// - create
	schedule, err := cronParser.Parse(spec)
	if err != nil {
		return nil, err
	}
	tr := &TimeRange{
		spec:     spec,
		schedule: schedule,
	}
	return tr, nil
}

func (r *TimeRange) IsTimeUpNow() (bool, time.Time) {
	return r.IsTimeUp(time.Now())
}

func (r *TimeRange) IsTimeUp(tm time.Time) (bool, time.Time) {
	next := r.schedule.Next(tm)
	if next.IsZero() {
		return false, next
	}

	// scheduler 最小单位是 min, 所以当相差小于 1 min时, 说明是当前可以运行的
	if next.Sub(tm) <= time.Minute {
		return true, tm
	}

	return false, next
}

func (r *TimeRange) Spec() string {
	return r.spec
}

func (t *TimeRange) Clock() (bool, time.Time) {
	return t.ClockIn(time.Now())
}

func (t *TimeRange) ClockIn(tm time.Time) (bool, time.Time) {
	return t.IsTimeUp(tm)
}
