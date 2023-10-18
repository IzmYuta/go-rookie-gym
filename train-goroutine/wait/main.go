package main

import (
	"log"
	"sync"
)

func main() {
	// waitgroupは内部のカウンタを用いるためfuncに渡すときはコピーではなく参照を渡す
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		i := i
		wg.Add(1)
		// ポインタを渡す
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			log.Println(i)
		}(&wg)// ポインタを渡す
	}
	wg.Wait()
	log.Println("end")
}
