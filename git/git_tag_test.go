package git

import (
	"reflect"
	"testing"
	"time"
)

func TestAnalyze(t *testing.T) {
	gt := &GitTag{
		path: "../glstats-sample-submodule",
	}

	t0 := time.Date(2017, time.Month(9), 1, 0, 0, 0, 0, time.Local)
	t1 := time.Date(2017, time.Month(10), 1, 0, 0, 0, 0, time.Local)
	t2 := time.Date(2017, time.Month(11), 1, 0, 0, 0, 0, time.Local)

	times := []time.Time{t0, t1, t2}

	expected := []TagCount{
		{time: t0, cnt: 0},
		{time: t1, cnt: 0},
		{time: t2, cnt: 1},
	}

	actual, err := gt.Analyze("deploy", times)
	if err != nil {
		t.Fatal("Analyze returnen non-nil\n")
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Actual was %x\n", actual)
	}
}
