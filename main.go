package main

import "runtime"

func main() {
	// This is the fastest permutation according to the benchmarks
	runtime.GOMAXPROCS(2)
	ch, _ := Producer(4, NewUUID)
	Consumer(2, ch)
}
