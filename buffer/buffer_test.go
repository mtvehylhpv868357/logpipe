package buffer

import (
	"sync"
	"testing"
	"time"
)

func collect() (func([]string), *[]string, *sync.Mutex) {
	var mu sync.Mutex
	var got []string
	fn := func(lines []string) {
		mu.Lock()
		got = append(got, lines...)
		mu.Unlock()
	}
	return fn, &got, &mu
}

func TestBuffer_FlushOnMaxSize(t *testing.T) {
	fn, got, mu := collect()
	b := New(3, 0, fn)
	b.Add("a")
	b.Add("b")
	b.Add("c") // should trigger flush
	mu.Lock()
	defer mu.Unlock()
	if len(*got) != 3 {
		t.Fatalf("expected 3 lines flushed, got %d", len(*got))
	}
}

func TestBuffer_ManualFlush(t *testing.T) {
	fn, got, mu := collect()
	b := New(100, 0, fn)
	b.Add("x")
	b.Add("y")
	b.Flush()
	mu.Lock()
	defer mu.Unlock()
	if len(*got) != 2 {
		t.Fatalf("expected 2, got %d", len(*got))
	}
}

func TestBuffer_EmptyFlushIsNoop(t *testing.T) {
	called := 0
	b := New(10, 0, func(lines []string) { called++ })
	b.Flush()
	if called != 0 {
		t.Fatalf("expected no flush call on empty buffer")
	}
}

func TestBuffer_TickerFlush(t *testing.T) {
	fn, got, mu := collect()
	b := New(100, 20*time.Millisecond, fn)
	b.Add("tick-line")
	time.Sleep(60 * time.Millisecond)
	b.Stop()
	mu.Lock()
	defer mu.Unlock()
	if len(*got) == 0 {
		t.Fatal("expected at least one line flushed via ticker")
	}
}

func TestBuffer_StopFlushesPending(t *testing.T) {
	fn, got, mu := collect()
	b := New(100, 0, fn)
	b.Add("pending")
	b.Stop()
	mu.Lock()
	defer mu.Unlock()
	if len(*got) != 1 {
		t.Fatalf("expected 1 line after Stop, got %d", len(*got))
	}
}

func TestBuffer_ConcurrentAdd(t *testing.T) {
	fn, got, mu := collect()
	b := New(5, 0, fn)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.Add("line")
		}()
	}
	wg.Wait()
	b.Flush()
	mu.Lock()
	defer mu.Unlock()
	if len(*got) != 50 {
		t.Fatalf("expected 50 total lines, got %d", len(*got))
	}
}
