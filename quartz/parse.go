package quartz

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Parse returns a new crontab schedule representing the given spec.
// 6个字段 用空格间隔
// (second) (minute) (hour) (day of month, optional) (month) (day of week, optional)
// ?只能用于dom和dow，表示不关心， 不能两个都不关心, 最终返回0
func Parse(spec string) (Schedule, error) {
	fields := strings.Fields(spec)
	if len(fields) != 6 {
		// log.Printf("Parse||Expected 6 fields, found %d: %s\n", len(fields), spec)
		return nil, fmt.Errorf("Expected 6 fields, found %d: %s", len(fields), spec)
	}

	for i := 0; i <= 2; i++ {
		if fields[i] == "?" {
			return nil, Error_WRONG_CRON
		}
	}
	if fields[4] == "?" {
		return nil, Error_WRONG_CRON
	}
	// 不允许dom和dow都无意义
	if fields[3] == "?" && fields[5] == "?" {
		return nil, Error_WRONG_CRON
	}
	// 不允许dom和dow都有意义
	if fields[3] != "?" && fields[5] != "?" {
		return nil, Error_WRONG_CRON
	}

	var err error
	var sec, min, hou, doom, mon, doow uint64
	sec, err = getField(fields[0], seconds)
	if err != nil {
		return nil, err
	}
	min, err = getField(fields[1], minutes)
	if err != nil {
		return nil, err
	}

	hou, err = getField(fields[2], hours)
	if err != nil {
		return nil, err
	}
	doom, err = getField(fields[3], dom)
	if err != nil {
		return nil, err
	}
	mon, err = getField(fields[4], months)
	if err != nil {
		return nil, err
	}
	doow, err = getField(fields[5], dow)
	if err != nil {
		return nil, err
	}
	schedule := &SpecSchedule{
		Second: sec,
		Minute: min,
		Hour:   hou,
		Dom:    doom,
		Month:  mon,
		Dow:    doow,
	}

	// schedule.printSpecSchedule()
	return schedule, nil
}

// getField returns an Int with the bits set representing all of the times that
// the field represents.  A "field" is a comma-separated list of "ranges".
// 每个字段可以由多个表达式组成，表达式间用逗号间隔
func getField(field string, r bounds) (uint64, error) {
	// list = range {"," range}
	var bits uint64
	ranges := strings.FieldsFunc(field, func(r rune) bool {
		return r == ','
	})
	for _, expr := range ranges {
		bit, err := getRange(expr, r)
		if err != nil {
			return 0, err
		}
		bits |= bit
	}
	return bits, nil
}

// 支持的表达式
// "*"
// "?"
// "5"
// "30/6"
// "2-10" []
func getRange(expr string, r bounds) (uint64, error) {

	var (
		start, end, step uint
		startAndStep     = strings.Split(expr, "/")
		lowAndHigh       = strings.Split(expr, "-")
		isStartAndStep   = (len(startAndStep) == 2)
		isLowAndHigh     = (len(lowAndHigh) == 2)
	)

	if isStartAndStep {
		tmp, err := strconv.Atoi(startAndStep[0])
		if tmp < 0 || err != nil {
			return 0, Error_WRONG_CRON
		}
		start = uint(tmp)
		tmp, err = strconv.Atoi(startAndStep[1])
		if tmp < 0 || err != nil {
			return 0, Error_WRONG_CRON
		}
		step = uint(tmp)
		end = r.max
	} else if isLowAndHigh {
		tmp, err := strconv.Atoi(lowAndHigh[0])
		if tmp < 0 || err != nil {
			return 0, Error_WRONG_CRON
		}
		start = uint(tmp)
		if start < r.min {
			return 0, Error_WRONG_CRON
		}
		step = 1
		tmp, err = strconv.Atoi(lowAndHigh[1])
		if tmp < 0 || err != nil {
			return 0, Error_WRONG_CRON
		}
		end = uint(tmp)
	} else {
		if expr == "*" {
			start = r.min
			end = r.max
			step = 1
		} else if expr == "?" {
			return 0, nil
		} else {
			tmp, err := strconv.Atoi(expr)
			if tmp < 0 || err != nil {
				return 0, Error_WRONG_CRON
			}
			start = uint(tmp)
			end = start
			step = 1
		}
	}

	if start < r.min || end > r.max || start > end {
		fmt.Println("wrong cron")
		return 0, Error_WRONG_CRON
	}

	return getBits(start, end, step), nil
}

func getBits(min, max, step uint) uint64 {
	var bits uint64

	// just a trick
	if step == 1 {
		return ^(math.MaxUint64 << (max + 1)) & (math.MaxUint64 << min)
	}

	// normal condition
	for i := min; i <= max; i += step {
		bits |= 1 << i
	}
	return bits
}
