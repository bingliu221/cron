package cron

import "time"

type IntSelector func(n int) bool

type TimeSelector struct {
	Second   IntSelector
	Minute   IntSelector
	Hour     IntSelector
	Day      IntSelector
	Month    IntSelector
	Weekday  IntSelector
	Location *time.Location
}

func SelectAll() IntSelector {
	return func(n int) bool {
		return true
	}
}

func SelectSlice(start int, end int, every int) IntSelector {
	return func(n int) bool {
		switch {
		case n < start:
			return false
		case n >= end:
			return false
		case (n-start)%every == 0:
			return true
		default:
			return false
		}
	}
}

func SelectSpecific(numbers ...int) IntSelector {
	m := make(map[int]bool)
	for _, n := range numbers {
		m[n] = true
	}
	return func(n int) bool {
		_, ok := m[n]
		return ok
	}
}
