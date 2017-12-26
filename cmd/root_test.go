package cmd

import (
	"testing"
	"time"
)

func TestGetByNow(t *testing.T) {
	expected := []time.Time{
		time.Date(2017, time.Month(3), 1, 0, 0, 0, 0, time.Local),
		time.Date(2017, time.Month(4), 1, 0, 0, 0, 0, time.Local),
		time.Date(2017, time.Month(5), 1, 0, 0, 0, 0, time.Local),
	}

	until := time.Date(2017, time.Month(6), 1, 0, 0, 0, 0, time.Local)
	dates := GetTimesUntil(until, 3, MONTH)

	if len(dates) != 3 {
		t.Errorf("dates length should be %d, but %t\n", expected, dates)
	}

	for i, v := range dates {
		if v != expected[i] {
			t.Errorf("date[%d] should be %x, but %x\n", i, expected[i], v)
		}
	}
}
