package main

import (
	"fmt"
	"sync"

	ulid "github.com/vvvvv/safeulid"
)

func main() {
	p := 8 // number of parallel coroutines
	number_of_ids := 50
	var mu sync.Mutex
	var all []ulid.ID

	var wg sync.WaitGroup
	for i := 0; i < p; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ids := make([]ulid.ID, number_of_ids)
			for j := 0; j < number_of_ids; j++ {
				// 0 * 100 + j 			= 	0
				//							1
				//							...
				//							99
				// 1 * 100 + j 			= 	100
				// 							101
				//							...
				ids[j] = ulid.MustNew()
			}
			mu.Lock()
			all = append(all, ids...)
			mu.Unlock()
		}()
	}
	wg.Wait()

	for _, id := range all {
		fmt.Printf("%v\n", id)
	}
	fmt.Println()
}
