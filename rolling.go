package rolling

import (
	"math"
	"time"

	"github.com/gammazero/deque"
)

type ll struct {
	value float64
	ts    time.Time
	next  *ll
	prev  *ll
}

// Trying different approaches to maintain min/max values:
// Linked list:
// BenchmarkWindow_Add_10k-8          55465             37098 ns/op
// BenchmarkWindow_Add_100k-8         65210             98093 ns/op
// RBT:
// BenchmarkWindow_Add_10k-8        1614670               631.1 ns/op
// BenchmarkWindow_Add_100k-8       1562439               709.3 ns/op
// Deque:
// BenchmarkWindow_Add_10k-8        5358396               200.0 ns/op
// BenchmarkWindow_Add_100k-8       4138747               315.4 ns/op

type Window struct {
	maxSize  int64
	duration time.Duration
	head     *ll
	tail     *ll

	cnt int64
	sum float64
	min float64
	max float64

	minDeque *deque.Deque[float64]
	maxDeque *deque.Deque[float64]
}

func NewWindow(maxSize int64, duration time.Duration) *Window {
	return &Window{
		maxSize:  maxSize,
		duration: duration,
		minDeque: deque.New[float64](),
		maxDeque: deque.New[float64](),
		min:      math.MaxFloat64,
		max:      -math.MaxFloat64,
	}
}

func (w *Window) addMinMax(value float64) {
	if value < w.min || w.cnt == 0 {
		w.min = value
		w.minDeque.Clear()
		w.minDeque.PushFront(value)
	} else {
		w.minDeque.PushBack(value)
	}
	if value > w.max || w.cnt == 0 {
		w.max = value
		w.minDeque.Clear()
		w.maxDeque.PushFront(value)
	} else {
		w.maxDeque.PushBack(value)
	}
}

func (w *Window) removeMinMax(value float64) {
	if value == w.max {
		w.maxDeque.PopFront()
		w.max = math.NaN()
		if w.maxDeque.Len() > 0 {
			w.max = w.maxDeque.Front()
		}
	}
	if value == w.min {
		w.minDeque.PopFront()
		w.min = math.NaN()
		if w.minDeque.Len() > 0 {
			w.min = w.minDeque.Front()
		}
	}
}

func (w *Window) Add(value float64) {
	w.cnt++
	w.sum += value

	if w.cnt > w.maxSize ||
		(w.head != nil && time.Since(w.head.ts) > w.duration) {
		w.sum -= w.head.value
		w.cnt--
		w.removeMinMax(w.head.value)
		w.head = w.head.next
	}

	w.addMinMax(value)

	if w.head == nil {
		w.head = &ll{value: value, ts: time.Now()}
		w.tail = w.head
		return
	}

	w.tail.next = &ll{
		value: value,
		ts:    time.Now(),
		prev:  w.tail,
	}
	w.tail = w.tail.next
}

func (w *Window) Sum() float64 {
	return w.sum
}

func (w *Window) Count() int64 {
	return w.cnt
}

func (w *Window) Min() float64 {
	return w.min
}

func (w *Window) Max() float64 {
	return w.max
}

func (w *Window) Avg() float64 {
	return w.sum / float64(w.cnt)
}

func (w *Window) Mid() float64 {
	return w.head.value + (w.tail.value-w.head.value)/2
}

func (w *Window) First() float64 {
	return w.head.value
}

func (w *Window) Last() float64 {
	return w.tail.value
}
