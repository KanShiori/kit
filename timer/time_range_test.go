package timer

import (
	"testing"
	"time"
)

func TestTimeRange(t *testing.T) {

	tr, err := NewHourTimeRange(0, 17)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(tr.Spec())

	tm, err := time.Parse("15:04:05", "19:01:12")
	if err != nil {
		t.Fatal(err)
	}

	is, next := tr.IsTimeUp(tm)
	t.Logf("is:%v next:%s\n", is, next)
}
