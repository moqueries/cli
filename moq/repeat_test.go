package moq_test

import (
	"testing"

	"github.com/myshkin5/moqueries/moq"
)

func TestRepeat(t *testing.T) {
	tests := map[string]struct {
		input        []moq.Repeater
		want         moq.RepeatVal
		fatalfFormat string
		fatalfArgs   []interface{}
	}{
		"default": {
			input:        nil,
			want:         moq.RepeatVal{MinTimes: 1, MaxTimes: 1, AnyTimes: false},
			fatalfFormat: "",
			fatalfArgs:   nil,
		},
		"min": {
			input:        []moq.Repeater{moq.MinTimes(2)},
			want:         moq.RepeatVal{MinTimes: 2, MaxTimes: 0, AnyTimes: true},
			fatalfFormat: "",
			fatalfArgs:   nil,
		},
		"max": {
			input:        []moq.Repeater{moq.MaxTimes(4)},
			want:         moq.RepeatVal{MinTimes: 0, MaxTimes: 4, AnyTimes: false},
			fatalfFormat: "",
			fatalfArgs:   nil,
		},
		"any": {
			input:        []moq.Repeater{moq.AnyTimes()},
			want:         moq.RepeatVal{MinTimes: 0, MaxTimes: 0, AnyTimes: true},
			fatalfFormat: "",
			fatalfArgs:   nil,
		},
		"times": {
			input:        []moq.Repeater{moq.Times(5)},
			want:         moq.RepeatVal{MinTimes: 5, MaxTimes: 5, AnyTimes: false},
			fatalfFormat: "",
			fatalfArgs:   nil,
		},
		"err/conflicting mins": {
			input:        []moq.Repeater{moq.MinTimes(6), moq.MinTimes(8)},
			want:         moq.RepeatVal{},
			fatalfFormat: "repeat min of %d conflicts with min of %d",
			fatalfArgs:   []interface{}{6, 8},
		},
		"err/conflicting maxs": {
			input:        []moq.Repeater{moq.MaxTimes(10), moq.MaxTimes(12)},
			want:         moq.RepeatVal{},
			fatalfFormat: "repeat max of %d conflicts with max of %d",
			fatalfArgs:   []interface{}{10, 12},
		},
		"err/conflicting times": {
			input:        []moq.Repeater{moq.Times(11), moq.Times(13)},
			want:         moq.RepeatVal{},
			fatalfFormat: "repeat times of %d conflicts with times of %d",
			fatalfArgs:   []interface{}{11, 13},
		},
		"err/max then any": {
			input:        []moq.Repeater{moq.MaxTimes(14), moq.AnyTimes()},
			want:         moq.RepeatVal{},
			fatalfFormat: "repeat max of %d conflicts with moq.AnyTimes",
			fatalfArgs:   []interface{}{14},
		},
		"err/any then max": {
			input:        []moq.Repeater{moq.AnyTimes(), moq.MaxTimes(16)},
			want:         moq.RepeatVal{},
			fatalfFormat: "repeat max of %d conflicts with moq.AnyTimes",
			fatalfArgs:   []interface{}{16},
		},
		"err/times then min": {
			input:        []moq.Repeater{moq.Times(14), moq.MinTimes(15)},
			want:         moq.RepeatVal{},
			fatalfFormat: "repeat min of %d conflicts with times %d",
			fatalfArgs:   []interface{}{15, 14},
		},
		"err/min then times": {
			input:        []moq.Repeater{moq.MinTimes(15), moq.Times(16)},
			want:         moq.RepeatVal{},
			fatalfFormat: "repeat min of %d conflicts with times %d",
			fatalfArgs:   []interface{}{15, 16},
		},
		"err/times then max": {
			input:        []moq.Repeater{moq.Times(14), moq.MaxTimes(15)},
			want:         moq.RepeatVal{},
			fatalfFormat: "repeat max of %d conflicts with times %d",
			fatalfArgs:   []interface{}{15, 14},
		},
		"err/max then times": {
			input:        []moq.Repeater{moq.MaxTimes(15), moq.Times(16)},
			want:         moq.RepeatVal{},
			fatalfFormat: "repeat max of %d conflicts with times %d",
			fatalfArgs:   []interface{}{15, 16},
		},
		"err/min then small max": {
			input:        []moq.Repeater{moq.MinTimes(18), moq.MaxTimes(9)},
			want:         moq.RepeatVal{},
			fatalfFormat: "repeat max of %d is less than min of %d",
			fatalfArgs:   []interface{}{9, 18},
		},
		"err/max then large min": {
			input:        []moq.Repeater{moq.MaxTimes(10), moq.MinTimes(20)},
			want:         moq.RepeatVal{},
			fatalfFormat: "repeat max of %d is less than min of %d",
			fatalfArgs:   []interface{}{10, 20},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// ASSEMBLE
			scene := moq.NewScene(t)
			tMoq := moq.NewMoqT(scene, nil)

			if test.fatalfFormat != "" {
				tMoq.OnCall().Fatalf(test.fatalfFormat, test.fatalfArgs...).ReturnResults()
			}

			// ACT
			got := moq.Repeat(tMoq.Mock(), test.input)

			// ASSERT
			scene.AssertExpectationsMet()

			if got != test.want {
				t.Errorf("Wanted %#v, got %#v", test.want, got)
			}
		})
	}
}
