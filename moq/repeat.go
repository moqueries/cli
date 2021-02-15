package moq

import (
	"github.com/myshkin5/moqueries/logs"
)

type Repeater interface {
	repeat()
}

type RepeatVal struct {
	MinTimes, MaxTimes int
	AnyTimes           bool
}

func Repeat(t T, repeaters []Repeater) RepeatVal {
	min, max, times := 1, 1, 1
	minDef, maxDef, timesDef, any := false, false, false, false
	for _, repeater := range repeaters {
		switch r := repeater.(type) {
		case minTimer:
			if minDef {
				t.Fatalf("repeat min of %d conflicts with min of %d", min, r)
				return RepeatVal{}
			}
			min = int(r)
			minDef = true
		case maxTimer:
			if maxDef {
				t.Fatalf("repeat max of %d conflicts with max of %d", max, r)
				return RepeatVal{}
			}
			max = int(r)
			maxDef = true
		case anyTimer:
			any = true
		case timer:
			if timesDef {
				t.Fatalf("repeat times of %d conflicts with times of %d", times, r)
				return RepeatVal{}
			}
			times = int(r)
			timesDef = true
		default:
			logs.Panicf("Unknown repeater type: %#v", repeater)
		}
	}

	if minDef && maxDef && min > max {
		t.Fatalf("repeat max of %d is less than min of %d", max, min)
		return RepeatVal{}
	}
	if maxDef && any {
		t.Fatalf("repeat max of %d conflicts with moq.AnyTimes", max)
		return RepeatVal{}
	}
	if minDef && timesDef {
		t.Fatalf("repeat min of %d conflicts with times %d", min, times)
		return RepeatVal{}
	}
	if maxDef && timesDef {
		t.Fatalf("repeat max of %d conflicts with times %d", max, times)
		return RepeatVal{}
	}

	if !minDef && maxDef || any {
		min = 0
	}
	if !maxDef && minDef || any {
		max = 0
		any = true
	}
	if timesDef {
		min = times
		max = times
	}

	return RepeatVal{
		MinTimes: min,
		MaxTimes: max,
		AnyTimes: any,
	}
}

type minTimer int

func (t minTimer) repeat() {}

func MinTimes(times int) minTimer {
	return minTimer(times)
}

type maxTimer int

func (t maxTimer) repeat() {}

func MaxTimes(times int) maxTimer {
	return maxTimer(times)
}

type anyTimer struct{}

func (t anyTimer) repeat() {}

func AnyTimes() anyTimer {
	return anyTimer{}
}

type timer int

func (t timer) repeat() {}

func Times(times int) timer {
	return timer(times)
}
