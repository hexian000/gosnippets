// gosnippets (c) 2023-2026 He Xian <hexian000@outlook.com>
// This code is licensed under MIT license (see LICENSE for details)

package algo

import (
	"math"
	"sync"
	"testing"
)

func TestSlidingDigest_Empty(t *testing.T) {
	d := NewSlidingDigest(100)
	for _, got := range []float64{d.Quantile(0), d.Quantile(0.5), d.Quantile(0.9), d.Quantile(0.99), d.Quantile(1)} {
		if !math.IsNaN(got) {
			t.Errorf("empty digest: got %v, want NaN", got)
		}
	}
}

func TestSlidingDigest_SingleValue(t *testing.T) {
	d := NewSlidingDigest(100)
	d.Add(42.0)
	cases := []struct {
		name string
		got  float64
	}{
		{"Min", d.Quantile(0)},
		{"P50", d.Quantile(0.5)},
		{"P90", d.Quantile(0.9)},
		{"P99", d.Quantile(0.99)},
		{"Max", d.Quantile(1)},
	}
	for _, c := range cases {
		if c.got != 42.0 {
			t.Errorf("%s = %v, want 42.0", c.name, c.got)
		}
	}
}

// TestSlidingDigest_Uniform verifies approximate quantile accuracy on a
// uniform distribution {1 … 100}.  Tolerance is ±2 units.
func TestSlidingDigest_Uniform(t *testing.T) {
	d := NewSlidingDigest(200)
	for i := 1; i <= 100; i++ {
		d.Add(float64(i))
	}
	const tol = 2.0
	cases := []struct {
		name string
		got  float64
		want float64
	}{
		{"P50", d.Quantile(0.5), 50.5},
		{"P90", d.Quantile(0.9), 90.5},
		{"P99", d.Quantile(0.99), 99.5},
	}
	for _, c := range cases {
		if math.Abs(c.got-c.want) > tol {
			t.Errorf("%s = %v, want %v ± %v", c.name, c.got, c.want, tol)
		}
	}
	if d.Quantile(0) != 1.0 {
		t.Errorf("Min = %v, want 1.0", d.Quantile(0))
	}
	if d.Quantile(1) != 100.0 {
		t.Errorf("Max = %v, want 100.0", d.Quantile(1))
	}
}

// TestSlidingDigest_WindowEviction verifies that observations older than the
// window size are discarded.
func TestSlidingDigest_WindowEviction(t *testing.T) {
	d := NewSlidingDigest(10)
	for i := 1; i <= 20; i++ {
		d.Add(float64(i))
	}
	// Window now holds [11 … 20].
	if d.Quantile(1) != 20.0 {
		t.Errorf("Max = %v, want 20.0", d.Quantile(1))
	}
	if d.Quantile(0) != 11.0 {
		t.Errorf("Min = %v, want 11.0", d.Quantile(0))
	}
	p50 := d.Quantile(0.5)
	// Median of {11 … 20} is 15.5 — exact for small windows (each value has
	// its own centroid when n << delta).
	if math.Abs(p50-15.5) > 1e-9 {
		t.Errorf("P50 = %v, want 15.5", p50)
	}
	// Confirm the first ten values (1…10) were evicted.
	if p50 < 11.0 {
		t.Errorf("P50 = %v, old values not evicted", p50)
	}
}

// TestSlidingDigest_NonFiniteInput verifies that NaN and ±Inf passed to Add
// are silently discarded and do not corrupt the digest.
func TestSlidingDigest_NonFiniteInput(t *testing.T) {
	d := NewSlidingDigest(10)
	d.Add(math.NaN())
	d.Add(math.Inf(1))
	d.Add(math.Inf(-1))
	// Non-finite values must not be counted.
	if got := d.Quantile(0.5); !math.IsNaN(got) {
		t.Errorf("after only non-finite Add: Quantile(0.5) = %v, want NaN", got)
	}
	// A subsequent finite value must be recorded correctly.
	d.Add(7.0)
	for _, q := range []float64{0, 0.5, 1} {
		if got := d.Quantile(q); got != 7.0 {
			t.Errorf("Quantile(%v) = %v, want 7.0", q, got)
		}
	}
	if got := d.Quantile(1); got != 7.0 {
		t.Errorf("Quantile(1) = %v, want 7.0", got)
	}
}

// TestSlidingDigest_NaNQuantile verifies that Quantile(NaN) returns NaN.
func TestSlidingDigest_NaNQuantile(t *testing.T) {
	d := NewSlidingDigest(10)
	d.Add(1.0)
	if got := d.Quantile(math.NaN()); !math.IsNaN(got) {
		t.Errorf("Quantile(NaN) = %v, want NaN", got)
	}
}

// TestSlidingDigest_Concurrent verifies that concurrent Add and Quantile calls
// do not race.  Run with -race.
func TestSlidingDigest_Concurrent(t *testing.T) {
	d := NewSlidingDigest(50)
	var wg sync.WaitGroup
	for g := 0; g < 4; g++ {
		wg.Add(1)
		go func(offset int) {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				d.Add(float64(offset*100 + i))
				_ = d.Quantile(0.5)
			}
		}(g)
	}
	wg.Wait()
}
