package main

import (
	"io/ioutil"
	"sync"

	"github.com/google/uuid"
)

const tenMillion = 10_000_000

func NewUUID() uuid.UUID { return uuid.New() }

func Producer(numWorkers int, uuidFunc func() uuid.UUID) (chan uuid.UUID, *sync.WaitGroup) {
	if numWorkers == 0 {
		panic("Number of workers == 0")
	}

	ch := make(chan uuid.UUID, tenMillion)
	var wg sync.WaitGroup

	idsPerGoroutine := tenMillion / numWorkers
	remainder := tenMillion % numWorkers
	f := func(numIds int) {
		defer wg.Done()
		for y := 0; y != numIds; y++ {
			ch <- uuidFunc()
		}
	}
	for w := 0; w != numWorkers; w++ {
		wg.Add(1)
		if w+1 != numWorkers {
			go f(idsPerGoroutine)
		} else {
			go f(idsPerGoroutine + remainder)
		}
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	return ch, &wg
}

func Consumer(numWorkers int, ch chan uuid.UUID) {
	var wg sync.WaitGroup

	for i := 0; i != numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range ch {
				ioutil.Discard.Write(i[:])
			}
		}()
	}
	wg.Wait()
}
