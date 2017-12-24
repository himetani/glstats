package timeutil

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const (
	layout string = "200601"
)

type DurationType int

const (
	DAY DurationType = iota
	MONTH
	YEAR
)

func Divide(s, u string, dur DurationType) ([]time.Time, error) {
	var since time.Time
	var until time.Time
	var dates []time.Time
	var err error

	if s == "" {
		since = time.Date(2014, 01, 01, 0, 0, 0, 0, time.Local)
	} else {
		since, err = getTime(s)
		if err != nil {
			return nil, errors.New("Error")
		}
	}

	if u == "" {
		until = time.Now()
	} else {
		until, err = getTime(u)
		if err != nil {
			return nil, errors.New("Error")
		}
	}

	i := 0
	for {
		var t time.Time
		switch dur {
		case MONTH:
			t = since.AddDate(0, i, 0)
			dates = append(dates, t)
		default:
		}

		if !t.Before(until) {
			break
		}

		i++
	}

	return dates, nil
}

func getTime(s string) (time.Time, error) {
	a := strings.Split(s, "-")

	if len(a) != 2 {
		return time.Now(), errors.New("-s option format is invalid")
	}

	y, _ := strconv.Atoi(a[0])
	m, _ := strconv.Atoi(a[1])

	return time.Date(y, time.Month(m), 1, 0, 0, 0, 0, time.Local), nil
}

func GetTimesUntil(until time.Time, num int, durationType DurationType) []time.Time {
	var since time.Time
	var times []time.Time

	switch durationType {
	case MONTH:
		since = until.AddDate(0, -num, 0)
	case DAY:
		since = until.AddDate(0, 0, -num)
	default:
		since = until
	}

	for i := 0; ; i++ {
		var t time.Time
		switch durationType {
		case MONTH:
			t = since.AddDate(0, i, 0)
		case DAY:
			since = since.AddDate(0, 0, i)
		default:
		}

		if t.After(until) || t.Equal(until) {
			break
		}
		times = append(times, t)
	}

	return times
}
