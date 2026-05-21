// gosnippets (c) 2023-2026 He Xian <hexian000@outlook.com>
// This code is licensed under MIT license (see LICENSE for details)

// Package algo provides common algorithm utilities.
package algo

import "sync"

// EWMA estimates the moving average of a stream of values using
// exponential weighting. The smoothing factor alpha (0 < alpha ≤ 1)
// determines how quickly older observations decay; alpha = 2/(N+1)
// approximates a sliding window of N samples.
//
// EWMA is safe for concurrent use.
type EWMA struct {
	mu          sync.Mutex
	alpha       float64
	value       float64
	initialized bool
}

// NewEWMA creates an EWMA estimator with the given smoothing factor
// alpha. Panics if alpha is not in the range (0, 1].
func NewEWMA(alpha float64) *EWMA {
	if alpha <= 0 || alpha > 1 {
		panic("algo: EWMA alpha must be in range (0, 1]")
	}
	return &EWMA{alpha: alpha}
}

// NewEWMAWindow creates an EWMA estimator that approximates a sliding
// window of n samples, using alpha = 2/(n+1). Panics if n < 1.
func NewEWMAWindow(n int) *EWMA {
	if n < 1 {
		panic("algo: EWMA window size must be >= 1")
	}
	return NewEWMA(2.0 / float64(n+1))
}

// Add incorporates a new observation into the estimate.
func (e *EWMA) Add(v float64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if !e.initialized {
		e.value = v
		e.initialized = true
		return
	}
	e.value = e.alpha*v + (1-e.alpha)*e.value
}

// Value returns the current EWMA estimate.
// Returns 0 if no observations have been added yet.
func (e *EWMA) Value() float64 {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.value
}
