package cron

import (
	"sort"
)

type Round struct {
	index    int
	values   []int
	reversed map[int]int
}

func NewRound(values []int, min, max int) *Round {
	r := &Round{
		index:    0,
		values:   make([]int, 0),
		reversed: make(map[int]int),
	}
	sort.Ints(values)
	for _, v := range values {
		if v < min || v > max {
			continue
		}
		_, ok := r.reversed[v]
		if !ok {
			r.values = append(r.values, v)
			r.reversed[v] = len(r.values) - 1
		}
	}
	return r
}

func (r *Round) Tick() bool {
	r.index = (r.index + 1) % len(r.values)
	return r.index == 0
}

func (r *Round) TickBack() bool {
	r.index = (r.index - 1) % len(r.values)
	return r.index == len(r.values)-1
}

func (r *Round) Value() int {
	return r.values[r.index]
}

func (r *Round) ShiftTo(v int) {
	i, ok := r.reversed[v]
	if ok {
		r.index = i
		return
	}

	for i = 0; i < len(r.values); i++ {
		if r.values[i] >= v {
			r.index = i
			return
		}
	}
}

func (r *Round) Contains(v int) bool {
	_, ok := r.reversed[v]
	return ok
}
