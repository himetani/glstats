package util

import (
	"testing"
	"time"
)

func TestGetTime(t *testing.T) {
	expected := time.Date(2010, time.Month(4), 1, 0, 0, 0, 0, time.Local)
	s, err := GetTime("2010-04")
	if err != nil {
		t.Error(err)
	}

	if !s.Equal(expected) {
		t.Errorf("Date parse from string to time is fail %s\n", s.String())
	}
}

func TestDivide(t *testing.T) {
	expected := []time.Time{
		time.Date(2017, time.Month(4), 1, 0, 0, 0, 0, time.Local),
		time.Date(2017, time.Month(5), 1, 0, 0, 0, 0, time.Local),
		time.Date(2017, time.Month(6), 1, 0, 0, 0, 0, time.Local),
	}

	dates, err := Divide("2017-04", "2017-06", MONTH)

	if err != nil {
		t.Errorf("%s\n", err.Error())
	}

	if len(dates) != 3 {
		t.Errorf("dates length should be %d, but %t\n", expected, dates)
	}

	for i, v := range dates {
		if v != expected[i] {
			t.Errorf("date[%d] should be %x, but %x\n", i, expected[i], v)
		}
	}
}
