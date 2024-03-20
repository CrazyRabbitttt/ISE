package util

import (
	"fmt"
	"github.com/thanhpk/randstr"
	"sync"
	"testing"
)

func TestRemovePunctuation(t *testing.T) {
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	mp := make(map[string]struct{})
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < 1000; j++ {
				uuid := randstr.String(10)
				mu.Lock()
				if j == 0 {
					fmt.Printf("generated uid:%s\n", uuid)
				}
				if _, ok := mp[uuid]; ok {
					fmt.Printf("error, id already exists....uuid is:%s\n", uuid)
				}
				mp[uuid] = struct{}{}
				mu.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Printf("The len of map:%d\n", len(mp))
}
