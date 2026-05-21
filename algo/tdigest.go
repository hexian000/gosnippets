// gosnippets (c) 2023-2026 He Xian <hexian000@outlook.com>
// This code is licensed under MIT license (see LICENSE for details)

// Package algo provides general algorithm utilities.
package algo

import (
	"math"
	"sort"
	"sync"
)

// centroid is a weighted point in a T-Digest: a cluster of values summarised
// by their mean and total weight.
type centroid struct {
	mean  float64
	count float64
}

// tdigest is a compact representation of a distribution using the T-Digest
// algorithm (Dunning, 2019) with the k2 scale function.
//
// The k2 scale function concentrates more centroids near the distribution
// extremes, yielding high accuracy for tail quantiles (P99, P1) at the cost
// of some accuracy near the median.
type tdigest struct {
	centroids []centroid
	total     float64
	delta     float64
}

// k2 computes the k2 scale-function value at quantile q given compression
// parameter delta:
//
//	k2(q) = (delta / 2π) × arcsin(2q − 1)
func k2(q, delta float64) float64 {
	return (delta / (2 * math.Pi)) * math.Asin(2*q-1)
}

// buildFromSorted constructs the digest from a pre-sorted slice of float64
// values using the batch T-Digest algorithm.  The slice may be empty, in
// which case the digest is reset to the empty state.
func (d *tdigest) buildFromSorted(vals []float64) {
	d.centroids = d.centroids[:0]
	n := len(vals)
	if n == 0 {
		d.total = 0
		return
	}
	d.total = float64(n)

	curMean := vals[0]
	curCount := 1.0
	cumBefore := 0.0

	for i := 1; i < n; i++ {
		// Quantile boundaries of the candidate merged centroid.
		qL := cumBefore / d.total
		qR := (cumBefore + curCount + 1) / d.total
		// Merge vals[i] into the current centroid when the k2 span stays
		// within one unit; otherwise finalise and start a new centroid.
		if k2(qR, d.delta)-k2(qL, d.delta) < 1.0 {
			// Welford online mean update.
			curCount++
			curMean += (vals[i] - curMean) / curCount
		} else {
			d.centroids = append(d.centroids, centroid{mean: curMean, count: curCount})
			cumBefore += curCount
			curMean = vals[i]
			curCount = 1.0
		}
	}
	d.centroids = append(d.centroids, centroid{mean: curMean, count: curCount})
}

// quantile estimates the value at quantile q (0 ≤ q ≤ 1).
// Returns math.NaN() when the digest is empty.
func (d *tdigest) quantile(q float64) float64 {
	if d.total == 0 || math.IsNaN(q) {
		return math.NaN()
	}
	n := len(d.centroids)
	if n == 1 || q <= 0 {
		return d.centroids[0].mean
	}
	if q >= 1 {
		return d.centroids[n-1].mean
	}

	// Each centroid c[i] is treated as a point mass at its mean, placed at
	// its cumulative midpoint: cum + c[i].count/2.  The target rank is
	// interpolated linearly between adjacent centroid midpoints.
	target := q * d.total
	cum := 0.0
	for i, c := range d.centroids {
		mid := cum + c.count/2
		if mid > target {
			if i == 0 {
				return c.mean
			}
			prevMid := cum - d.centroids[i-1].count/2
			t := (target - prevMid) / (mid - prevMid)
			return d.centroids[i-1].mean + t*(c.mean-d.centroids[i-1].mean)
		}
		cum += c.count
	}
	return d.centroids[n-1].mean
}

// SlidingDigest estimates quantile statistics over the most recent N
// observations using the T-Digest algorithm (Dunning, 2019).
//
// The digest uses the k2 scale function with compression parameter δ = 100,
// giving sub-percent error for P99 and exact values near the extremes for
// small windows.  The maximum value is always exact.
//
// The digest is rebuilt lazily from the ring buffer when a query method is
// called after one or more Add calls.  Rebuild cost is O(N log N) where N is
// the window size.
//
// All methods are safe for concurrent use.
type SlidingDigest struct {
	mu      sync.Mutex
	buf     []float64
	head    int
	size    int
	dirty   bool
	min     float64
	max     float64
	scratch []float64
	td      tdigest
}

// NewSlidingDigest returns a SlidingDigest that retains the last windowSize
// observations.  windowSize must be at least 1.
func NewSlidingDigest(windowSize int) *SlidingDigest {
	if windowSize < 1 {
		panic("algo: SlidingDigest windowSize must be at least 1")
	}
	return &SlidingDigest{
		buf:     make([]float64, windowSize),
		scratch: make([]float64, windowSize),
		min:     math.NaN(),
		max:     math.NaN(),
		td:      tdigest{delta: 100.0},
	}
}

// Add records a new observation, evicting the oldest when the window is full.
// Non-finite values (NaN, ±Inf) are silently ignored.
func (d *SlidingDigest) Add(v float64) {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return
	}
	d.mu.Lock()
	d.buf[d.head] = v
	d.head = (d.head + 1) % len(d.buf)
	if d.size < len(d.buf) {
		d.size++
	}
	d.dirty = true
	d.mu.Unlock()
}

// rebuild recomputes the T-Digest from the ring buffer contents.
// Must be called with d.mu held.
func (d *SlidingDigest) rebuild() {
	if !d.dirty {
		return
	}
	n := d.size
	if n == 0 {
		d.td.total = 0
		d.td.centroids = d.td.centroids[:0]
		d.min = math.NaN()
		d.max = math.NaN()
		d.dirty = false
		return
	}
	tmp := d.scratch[:n]
	start := (d.head - n + len(d.buf)) % len(d.buf)
	minVal, maxVal := math.Inf(1), math.Inf(-1)
	for i := 0; i < n; i++ {
		v := d.buf[(start+i)%len(d.buf)]
		tmp[i] = v
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}
	sort.Float64s(tmp)
	d.td.buildFromSorted(tmp)
	d.min = minVal
	d.max = maxVal
	d.dirty = false
}

// Quantile estimates the value at quantile q.
//
// For q ≤ 0 or q ≥ 1 the exact minimum or maximum of the window is returned.
// Returns math.NaN() if no observations have been recorded or q is NaN.
func (d *SlidingDigest) Quantile(q float64) float64 {
	if math.IsNaN(q) {
		return math.NaN()
	}
	d.mu.Lock()
	d.rebuild()
	var v float64
	switch {
	case d.td.total == 0:
		v = math.NaN()
	case q <= 0:
		v = d.min
	case q >= 1:
		v = d.max
	default:
		v = d.td.quantile(q)
	}
	d.mu.Unlock()
	return v
}
