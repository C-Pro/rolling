package rolling

import (
	"math"
	"time"

	rbt "github.com/ugurcsen/gods-generic/trees/redblacktree"
	"github.com/ugurcsen/gods-generic/utils"
)

type ll struct {
	value float64
	ts    time.Time
	next  *ll
	prev  *ll
}

// Using RBT instead of linked list for distinct values.
// Linked list:
// BenchmarkWindow_Add_10k-8          55465             37098 ns/op
// BenchmarkWindow_Add_100k-8         65210             98093 ns/op
// RBT:
// BenchmarkWindow_Add_10k-8        1614670               631.1 ns/op
// BenchmarkWindow_Add_100k-8       1562439               709.3 ns/op
type distinctVal struct {
	value float64
	cnt   int64
}

type Window struct {
	maxSize  int64
	duration time.Duration
	head     *ll
	tail     *ll

	cnt int64
	sum float64
	min float64
	max float64

	orderedIndex *rbt.Tree[float64, *distinctVal]
}

func NewWindow(maxSize int64, duration time.Duration) *Window {
	return &Window{
		maxSize:      maxSize,
		duration:     duration,
		orderedIndex: rbt.NewWith[float64, *distinctVal](utils.NumberComparator[float64]),
	}
}

func (w *Window) addMinMax(value float64) {
	if value < w.min || w.cnt == 0 {
		w.min = value
	}
	if value > w.max || w.cnt == 0 {
		w.max = value
	}

	val, ok := w.orderedIndex.Get(value)
	if !ok {
		val = &distinctVal{value: value, cnt: 0}
	}

	val.cnt++
	w.orderedIndex.Put(value, val)
}

func (w *Window) removeMinMax(value float64) {
	val, _ := w.orderedIndex.Get(value)
	val.cnt--
	if val.cnt > 0 {
		return
	}

	if w.max == value {
		max := math.NaN()
		node := w.orderedIndex.GetNode(value)
		if node != nil {
			it := w.orderedIndex.IteratorAt(node)
			if it.Prev() {
				max = it.Value().value
			}
		}
		w.max = max
	}

	if w.min == value {
		min := math.NaN()
		node := w.orderedIndex.GetNode(value)
		if node != nil {
			it := w.orderedIndex.IteratorAt(node)
			if it.Next() {
				min = it.Value().value
			}
		}
		w.min = min
	}

	w.orderedIndex.Remove(value)
}

func (w *Window) Add(value float64) {
	w.addMinMax(value)
	w.cnt++
	w.sum += value

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

	if w.cnt > w.maxSize ||
		(w.head != nil && time.Since(w.head.ts) > w.duration) {
		w.sum -= w.head.value
		w.removeMinMax(w.head.value)
		w.head = w.head.next
		w.cnt--
	}
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
