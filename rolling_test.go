package rolling

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

func ExampleNewWindow() {
	w := NewWindow(6, time.Second)
	w.Add(1)
	w.Add(2)
	time.Sleep(time.Second)
	w.Add(3)
	w.Add(4)
	w.Add(5)
	w.Add(6)
	w.Add(7)
	fmt.Println(w.Min(), w.Avg(), w.Max())
	// Output: 3 5 7
}

const tolerance = 0.00000000001

func TestMin(t *testing.T) {
	tests := []struct {
		name   string
		wsize  int
		values []float64
		expect float64
	}{
		{
			"zero values",
			3,
			nil,
			math.NaN(),
		},
		{
			"1 value",
			3,
			[]float64{42},
			float64(42),
		},
		{
			"first in window",
			3,
			[]float64{1, 2, 3},
			float64(1),
		},
		{
			"middle in window",
			3,
			[]float64{2, 1, 3},
			float64(1),
		},
		{
			"last in window",
			3,
			[]float64{2, 3, 1},
			float64(1),
		},
		{
			"last in window, evict",
			3,
			[]float64{1, 3, 4, 2},
			float64(2),
		},
		{
			"first in window, evict same",
			3,
			[]float64{1, 1, 4, 2},
			float64(1),
		},
		{
			"first in window, evict",
			3,
			[]float64{1, 2, 4, 3},
			float64(2),
		},
		{
			"middle in window, evict",
			3,
			[]float64{1, 3, 2, 5},
			float64(2),
		},
	}

	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			w := NewWindow(int64(tcase.wsize), time.Hour)
			for _, v := range tcase.values {
				w.Add(v)
			}
			got := w.Min()
			if math.Abs(got-tcase.expect) > tolerance {
				t.Errorf("expected min %f, but got %f", tcase.expect, got)
			}
		})
	}
}

func arrayMax(a []float64) float64 {
	max := -math.MaxFloat64
	for _, v := range a {
		if v > max {
			max = v
		}
	}
	return max
}

func arrayMin(a []float64) float64 {
	min := math.MaxFloat64
	for _, v := range a {
		if v < min {
			min = v
		}
	}
	return min
}

func arrayAvg(a []float64) float64 {
	sum := float64(0)
	for _, v := range a {
		sum += v
	}
	return sum / float64(len(a))
}

func TestFuzz(t *testing.T) {
	const (
		sz         = 5
		iterations = 100_000
	)
	w := NewWindow(sz, time.Minute)
	window := make([]float64, sz)
	for i := 0; i < iterations; i++ {
		value := rand.Float64()
		w.Add(value)
		window[i%sz] = value
		if i > sz {
			avg := w.Avg()
			exp := arrayAvg(window)
			if math.Abs(avg-exp) > tolerance {
				t.Fatalf("%d: expected avg %f, got %f", i, exp, avg)
			}
			min := w.Min()
			exp = arrayMin(window)
			if math.Abs(min-exp) > tolerance {
				t.Fatalf("%d: expected min %f, got %f", i, exp, min)
			}
			max := w.Max()
			exp = arrayMax(window)
			if math.Abs(max-exp) > tolerance {
				t.Fatalf("%d: expected max %f, got %f", i, exp, max)
			}
		}
	}
}

func BenchmarkWindow_Add_10k(b *testing.B) {
	w := NewWindow(10_000, time.Second)
	for i := 0; i < b.N; i++ {
		w.Add(rand.Float64())
	}
}

func BenchmarkWindow_Add_100k(b *testing.B) {
	w := NewWindow(100_000, time.Second)
	for i := 0; i < b.N; i++ {
		w.Add(rand.Float64())
	}
}
