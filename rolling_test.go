package rolling

import (
	"fmt"
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

func BenchmarkWindow_Add_10k(b *testing.B) {
	w := NewWindow(10_000, time.Second)
	for i := 0; i < b.N; i++ {
		w.Add(float64(i))
	}
}

func BenchmarkWindow_Add_100k(b *testing.B) {
	w := NewWindow(100_000, time.Second)
	for i := 0; i < b.N; i++ {
		w.Add(float64(i))
	}
}
