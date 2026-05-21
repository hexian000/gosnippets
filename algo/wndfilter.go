// gosnippets (c) 2023-2026 He Xian <hexian000@outlook.com>
// This code is licensed under MIT license (see LICENSE for details)

// Package algo provides generic algorithm utilities.
package algo

// ordered is a constraint for types supporting comparison operators.
type ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 | ~string
}

type sample[V any] struct {
	t int64
	v V
}

// WindowedFilter implements the 3-slot windowed min/max filter described by
// Kathleen Nichols (2012). It tracks the best (minimum or maximum) value
// observed within a sliding time window using constant space and O(1) updates.
//
// V is the sample value type.
// Use [NewMinFilter] or [NewMaxFilter] to construct a filter.
type WindowedFilter[V ordered] struct {
	win         int64
	notWorse    func(a, b V) bool
	s           [3]sample[V]
	initialized bool
}

// NewMinFilter returns a WindowedFilter that tracks the minimum value within
// a sliding window of duration win.
func NewMinFilter[V ordered](win int64) *WindowedFilter[V] {
	return &WindowedFilter[V]{
		win:      win,
		notWorse: func(a, b V) bool { return a <= b },
	}
}

// NewMaxFilter returns a WindowedFilter that tracks the maximum value within
// a sliding window of duration win.
func NewMaxFilter[V ordered](win int64) *WindowedFilter[V] {
	return &WindowedFilter[V]{
		win:      win,
		notWorse: func(a, b V) bool { return a >= b },
	}
}

// reset restores the filter to a single-sample state at time t with value v.
func (f *WindowedFilter[V]) reset(t int64, v V) V {
	s := sample[V]{t: t, v: v}
	f.s[0] = s
	f.s[1] = s
	f.s[2] = s
	f.initialized = true
	return v
}

// subwinUpdate advances the subwindow slots as time progresses, ensuring the
// three samples remain spread across the window. Called after the value-based
// slot updates in Update.
func (f *WindowedFilter[V]) subwinUpdate(t int64, v V) V {
	dt := t - f.s[0].t
	val := sample[V]{t: t, v: v}

	if dt > f.win {
		// The best sample has aged out. Promote: 1st←2nd, 2nd←3rd, 3rd←new.
		// Repeat once because the promoted 2nd may also have aged out.
		f.s[0] = f.s[1]
		f.s[1] = f.s[2]
		f.s[2] = val
		if t-f.s[0].t > f.win {
			f.s[0] = f.s[1]
			f.s[1] = f.s[2]
			f.s[2] = val
		}
	} else if f.s[1].t == f.s[0].t && dt > f.win/4 {
		// A quarter-window has passed without a distinct 2nd sample.
		f.s[1] = val
		f.s[2] = val
	} else if f.s[2].t == f.s[1].t && dt > f.win/2 {
		// A half-window has passed without a distinct 3rd sample.
		f.s[2] = val
	}
	return f.s[0].v
}

// Update records a new sample with value v observed at time t, then returns
// the current best (min or max) value within the window. t must be
// monotonically non-decreasing across successive calls.
func (f *WindowedFilter[V]) Update(t int64, v V) V {
	// Reject NaN: IEEE 754 NaN is unordered, so all comparisons with it
	// return false. A NaN sample would bypass the value checks below and
	// corrupt slot values via the time-based logic in subwinUpdate.
	// For non-float V this comparison is always false and is eliminated
	// by the compiler.
	if v != v {
		return f.s[0].v
	}
	// Reset when uninitialized, a new best is found, or the window is empty.
	if !f.initialized || f.notWorse(v, f.s[0].v) || t-f.s[2].t > f.win {
		return f.reset(t, v)
	}
	// Update the 2nd and/or 3rd slots if the new value qualifies.
	if f.notWorse(v, f.s[1].v) {
		f.s[2] = sample[V]{t: t, v: v}
		f.s[1] = f.s[2]
	} else if f.notWorse(v, f.s[2].v) {
		f.s[2] = sample[V]{t: t, v: v}
	}
	return f.subwinUpdate(t, v)
}

// Get returns the current best (min or max) value without updating the filter.
// Returns the zero value of V if the filter has not been updated yet.
func (f *WindowedFilter[V]) Get() V {
	return f.s[0].v
}
