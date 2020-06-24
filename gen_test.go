package main

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/google/uuid"
)

func cheapNewUUID() uuid.UUID { var b [16]byte; return b }

func TestProducer(t *testing.T) {
	// Test with 1, 2, and 3 workers
	for i := 1; i != 4; i++ {
		ch, wg := Producer(i, cheapNewUUID)
		wg.Wait()
		if len(ch) != tenMillion {
			t.Errorf("Expected tenMillion, got %v", len(ch))
		}
	}
}

func BenchmarkNewUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewUUID()
	}
}

func BenchmarkNewUUIDTenMillion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for x := 0; x != tenMillion; x++ {
			NewUUID()
		}
	}
}

func BenchmarkProducer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, wg := Producer(2, NewUUID)
		wg.Wait()
	}
}

func BenchmarkConsumer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ch, wg := Producer(2, cheapNewUUID)
		wg.Wait()
		b.StartTimer()
		Consumer(2, ch)
	}
}

func BenchmarkRun(b *testing.B) {
	type benchmark struct {
		name      string
		producers int
		consumers int
	}
	var benchmarks []benchmark
	for p := 1; p != runtime.NumCPU()+1; p++ {
		for c := 1; c != runtime.NumCPU()+1; c++ {
			benchmarks = append(benchmarks, benchmark{fmt.Sprintf("%v-%v", p, c), p, c})
		}
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ch, _ := Producer(bm.producers, NewUUID)
				Consumer(bm.consumers, ch)
			}
		})
	}
}
