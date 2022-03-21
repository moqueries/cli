package moq

import (
	"github.com/myshkin5/moqueries/logs"
)

// Repeater is implemented by all repeaters
type Repeater interface {
	repeat()
}

// RepeatVal is a compilation of multiple repeaters for use in a moq
type RepeatVal struct {
	// MinTimes and MaxTimes record the minimum and maximum number of times a
	// call must be made for a test to pass
	MinTimes, MaxTimes int
	// AnyTimes indicates that a call can be made any number of times
	AnyTimes bool
	// ResultCount is used by a moq to size the results slice large enough to
	// hold all the expected results
	ResultCount int
	// ExplicitAny indicates that AnyTimes is true from an explicit use of
	// moq.AnyTimes
	ExplicitAny bool
	// Incremented indicates that Increment has been called as required before
	// calling Repeat
	Incremented bool
}

// Increment must be called by a moq each time results are expected and must be
// called prior to the optional call to Repeat
func (v *RepeatVal) Increment(t T) {
	t.Helper()
	v.repeat(t, []Repeater{Times(1)})
	v.Incremented = true
}

// Repeat compiles multiple repeaters into a value that can be used by a moq.
// If Repeat detects any rule violations (e.g.: conflicting repeaters),
// t.Fatalf is called to stop the test.
func (v *RepeatVal) Repeat(t T, repeaters []Repeater) {
	t.Helper()
	if !v.Incremented {
		t.Fatalf("fn Increment not called before fn Repeat")
		return
	}

	v.MinTimes--
	v.MaxTimes--
	v.Incremented = false

	v.repeat(t, repeaters)
}

func (v *RepeatVal) repeat(t T, repeaters []Repeater) {
	min, max, times := 1, 1, 1
	minDef, maxDef, timesDef, any, opt, explicitAny := false, false, false, false, false, false
	for _, repeater := range repeaters {
		switch r := repeater.(type) {
		case MinTimer:
			if minDef {
				t.Fatalf("repeat min of %d conflicts with min of %d", min, r)
				return
			}
			min = int(r)
			minDef = true
		case MaxTimer:
			if maxDef {
				t.Fatalf("repeat max of %d conflicts with max of %d", max, r)
				return
			}
			max = int(r)
			maxDef = true
		case AnyTimer:
			any = true
			explicitAny = true
		case OptionalTimer:
			opt = true
		case Timer:
			if timesDef {
				t.Fatalf("repeat times of %d conflicts with times of %d", times, r)
				return
			}
			times = int(r)
			timesDef = true
		default:
			logs.Panicf("Unknown repeater type: %#v", repeater)
		}
	}

	if !minDef && (maxDef || any || opt) {
		min = 0
	}
	if !maxDef && (minDef || any) && !opt {
		max = 0
		any = true
	}

	if minDef && timesDef {
		t.Fatalf("min of %d conflicts with times %d", min, times)
		return
	}
	if maxDef && timesDef {
		t.Fatalf("max of %d conflicts with times %d", max, times)
		return
	}
	if maxDef && any {
		t.Fatalf("max of %d conflicts with moq.AnyTimes", max)
		return
	}
	if minDef && opt {
		t.Fatalf("min of %d conflicts with moq.Optional", min)
		return
	}
	if timesDef {
		if !opt {
			min = times
		}
		max = times
	}

	v.MinTimes += min
	v.MaxTimes += max
	v.ExplicitAny = explicitAny || v.ExplicitAny
	if v.ExplicitAny {
		v.AnyTimes = true
	} else if v.MaxTimes == 0 {
		v.AnyTimes = any || v.AnyTimes
	}
	v.ResultCount = v.MinTimes
	if v.MaxTimes > v.ResultCount {
		v.ResultCount = v.MaxTimes
	}
	if v.AnyTimes {
		// Always one extra result for anyTimes
		v.ResultCount++
	}

	if v.MinTimes > v.MaxTimes && v.MaxTimes > 0 && !v.AnyTimes {
		t.Fatalf("max of %d is less than min of %d", v.MaxTimes, v.MinTimes)
		return
	}
}

// MinTimer holds a minimum times value
type MinTimer int

func (t MinTimer) repeat() {}

// MinTimes returns a minimum times value
func MinTimes(times int) MinTimer {
	return MinTimer(times)
}

// MaxTimer holds a maximum times value
type MaxTimer int

func (t MaxTimer) repeat() {}

// MaxTimes returns a maximum times value
func MaxTimes(times int) MaxTimer {
	return MaxTimer(times)
}

// AnyTimer holds an any times value
type AnyTimer struct{}

func (t AnyTimer) repeat() {}

// AnyTimes returns an any times value
func AnyTimes() AnyTimer {
	return AnyTimer{}
}

// OptionalTimer holds an optional value
type OptionalTimer struct{}

func (t OptionalTimer) repeat() {}

// Optional is similar to calling MinTimes(0) but doesn't change the max
func Optional() OptionalTimer {
	return OptionalTimer{}
}

// Timer holds a times value (min and max set to same value)
type Timer int

func (t Timer) repeat() {}

// Times returns a times value (min and max set to same value)
func Times(times int) Timer {
	return Timer(times)
}
