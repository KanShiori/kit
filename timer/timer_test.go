package timer

import (
	"testing"
	"time"
)

func TestTimerList(t *testing.T) {
	tr, err := NewHourTimeRange(0, 18)
	if err != nil {
		t.Fatal(err)
	}

	ts := NewTimeSpan(time.Hour, time.Time{})

	timer, err := NewTimerList(tr, ts)
	if err != nil {
		t.Fatal(err)
	}

	tm, err := time.Parse("2006 15:04:05", "2006 17:01:12")
	if err != nil {
		t.Fatal(err)
	}
	is, next := timer.ClockIn(tm)
	t.Logf("is:%v next:%s\n", is, next)

	tm, err = time.Parse("2006 15:04:05", "2006 20:01:12")
	if err != nil {
		t.Fatal(err)
	}
	is, next = timer.ClockIn(tm)
	t.Logf("is:%v next:%s\n", is, next)

	tm, err = time.Parse("2006 15:04:05", "2006 17:30:12")
	if err != nil {
		t.Fatal(err)
	}
	is, next = timer.ClockIn(tm)
	t.Logf("is:%v next:%s\n", is, next)

	tm, err = time.Parse("2006 15:04:05", "2006 18:01:12")
	if err != nil {
		t.Fatal(err)
	}
	is, next = timer.ClockIn(tm)
	t.Logf("is:%v next:%s\n", is, next)

	tm, err = time.Parse("2006 15:04:05", "2006 20:01:12")
	if err != nil {
		t.Fatal(err)
	}
	is, next = timer.ClockIn(tm)
	t.Logf("is:%v next:%s\n", is, next)
}
