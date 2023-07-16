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
// BenchmarkWindow_Add_10k-8        5558586               213.1 ns/op
// BenchmarkWindow_Add_100k-8       5055327               211.9 ns/op

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
	if w.minDeque.Len() > 0 {
		for w.minDeque.Len() > 0 {
			b := w.minDeque.Back()
			if value < b {
				w.minDeque.PopBack()
			} else {
				break
			}
		}
	}
	w.minDeque.PushBack(value)
	w.min = w.minDeque.Front()

	if w.maxDeque.Len() > 0 {
		for w.maxDeque.Len() > 0 {
			b := w.maxDeque.Back()
			if value > b {
				w.maxDeque.PopBack()
			} else {
				break
			}
		}
	}
	w.maxDeque.PushBack(value)
	w.max = w.maxDeque.Front()
}

func (w *Window) removeMinMax(value float64) {
	if w.maxDeque.Front() == value {
		w.maxDeque.PopFront()
	}
	if w.maxDeque.Len() == 0 {
		w.max = -math.MaxFloat64
	}
	if w.minDeque.Front() == value {
		w.minDeque.PopFront()
	}
	if w.minDeque.Len() == 0 {
		w.min = math.MaxFloat64
	}
}

func (w *Window) Add(value float64) {
	w.cnt++
	w.sum += value

	// Remove head if window is full.
	if w.cnt > w.maxSize {
		w.sum -= w.head.value
		w.cnt--
		w.removeMinMax(w.head.value)
		w.head = w.head.next
	}

	// Truncate old values.
	for w.head != nil && time.Since(w.head.ts) > w.duration {
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
