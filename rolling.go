package rolling

import "time"

type ll struct {
	value float64
	ts    time.Time
	next  *ll
	prev  *ll
}

// Linked list used to keep track of distinct values to support
// updating Min/Max without scanning the entire window.
type setLL struct {
	value float64
	cnt   int64
	next  *setLL
	prev  *setLL
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

	orderedHead  *setLL
	orderedTail  *setLL
	orderedIndex map[float64]*setLL
}

func NewWindow(maxSize int64, duration time.Duration) *Window {
	return &Window{
		maxSize:      maxSize,
		duration:     duration,
		orderedIndex: make(map[float64]*setLL),
	}
}

func (w *Window) addMinMax(value float64) {
	if value < w.min || w.cnt == 0 {
		w.min = value
	}
	if value > w.max || w.cnt == 0 {
		w.max = value
	}

	val, ok := w.orderedIndex[value]
	if !ok {
		val = &setLL{value: value, cnt: 0}
	}

	val.cnt++
	w.orderedIndex[value] = val

	if ok {
		// The value is already in the list, we're done here.
		return
	}

	if w.orderedHead == nil {
		w.orderedHead = val
		w.orderedTail = val
		return
	}

	if value < w.orderedHead.value {
		w.orderedHead.prev = val
		val.next = w.orderedHead
		w.orderedHead = val
		return
	}

	if value > w.orderedTail.value {
		w.orderedTail.next = val
		val.prev = w.orderedTail
		w.orderedTail = val
		return
	}

	// Search for the correct position to insert the new value.
	for cur := w.orderedHead; cur != nil; cur = cur.next {
		if value < cur.value {
			cur.prev.next = val
			val.prev = cur.prev
			val.next = cur
			cur.prev = val
			return
		}
	}
}

func (w *Window) removeMinMax(value float64) {
	val := w.orderedIndex[value]
	val.cnt--
	if val.cnt > 0 {
		return
	}

	delete(w.orderedIndex, value)

	if val == w.orderedHead {
		w.orderedHead = val.next
		w.orderedHead.prev = nil
		w.min = w.orderedHead.value
		return
	}

	if val == w.orderedTail {
		w.orderedTail = val.prev
		w.orderedTail.next = nil
		w.max = w.orderedTail.value
		return
	}

	val.prev.next = val.next
	val.next.prev = val.prev
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
