package main

import (
	"sync"

	"github.com/goodgoodjm/k-pioneer/batch"
)

func main() {
	batch.Start()
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
