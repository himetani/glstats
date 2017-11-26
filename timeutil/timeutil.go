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

type Duration int

const (
	MONTH Duration = iota
	YEAR
)

func Divide(s, u string, dur Duration) ([]time.Time, error) {
	var since time.Time
	var until time.Time
	var dates []time.Time
	var err error

	if s == "" {
		since = time.Date(2014, 01, 01, 0, 0, 0, 0, time.Local)
	} else {
		since, err = GetTime(s)
		if err != nil {
			return nil, errors.New("Error")
		}
	}

	if u == "" {
		until = time.Now()
	} else {
		until, err = GetTime(u)
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

func GetTime(s string) (time.Time, error) {
	a := strings.Split(s, "-")

	if len(a) != 2 {
		return time.Now(), errors.New("-s option format is invalid")
	}

	y, _ := strconv.Atoi(a[0])
	m, _ := strconv.Atoi(a[1])

	return time.Date(y, time.Month(m), 1, 0, 0, 0, 0, time.Local), nil
}
