package quartz

import (
	"fmt"
	"time"
)

// SpecSchedule specifies a duty cycle (to the second granularity), based on a
// traditional crontab specification. It is computed initially and stored as bit sets.
type SpecSchedule struct {
	Second, Minute, Hour, Dom, Month, Dow uint64
}

func (ss *SpecSchedule) printSpecSchedule() {
	fmt.Printf("Second:%b\n", ss.Second)
	fmt.Printf("Minute:%b\n", ss.Minute)
	fmt.Printf("Hour:%b\n", ss.Hour)
	fmt.Printf("Dom:%b\n", ss.Dom)
	fmt.Printf("Month:%b\n", ss.Month)
	fmt.Printf("Dow:%b\n", ss.Dow)

}

// bounds provides a range of acceptable values (plus a map of name to value).
type bounds struct {
	min, max uint
}

// The bounds for each field.
var (
	seconds = bounds{0, 59}
	minutes = bounds{0, 59}
	hours   = bounds{0, 23}
	dom     = bounds{1, 31}
	months  = bounds{1, 12}
	dow     = bounds{0, 6}
)

func (s *SpecSchedule) Next(t time.Time) time.Time {

	t = t.Add(1*time.Second - time.Duration(t.Nanosecond())*time.Nanosecond)
	added := false
	yearLimit := t.Year() + 5

WRAP:
	if t.Year() > yearLimit {
		return time.Time{}
	}

	for (1<<uint(t.Month()))&s.Month == 0 {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
		}
		t = t.AddDate(0, 1, 0)
		if t.Month() == time.January {
			goto WRAP
		}
	}

	if s.Dom != 0 {
		for (1<<uint(t.Day()))&s.Dom == 0 {
			if !added {
				added = true
				t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
			}
			t = t.AddDate(0, 0, 1)
			if t.Day() == 1 {
				goto WRAP
			}
		}
	} else {
		for (1<<uint(t.Weekday()))&s.Dow == 0 {
			if !added {
				added = true
				t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
			}
			t = t.AddDate(0, 0, 1)
			if t.Day() == 1 {
				goto WRAP
			}
		}
	}
	for (1<<uint(t.Hour()))&s.Hour == 0 {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
		}
		t = t.Add(1 * time.Hour)
		if t.Hour() == 0 {
			goto WRAP
		}
	}

	for (1<<uint(t.Minute()))&s.Minute == 0 {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
		}
		t = t.Add(1 * time.Minute)
		if t.Minute() == 0 {
			goto WRAP
		}
	}

	for (1<<uint(t.Second()))&s.Second == 0 {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
		}
		t = t.Add(1 * time.Second)

		if t.Second() == 0 {
			goto WRAP
		}
	}

	return t
}
