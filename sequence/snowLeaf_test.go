package sequence

import (
	"Search-Engine/search-engine/util"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"sync"
	"testing"
)

func TestBytes2IntFunction(t *testing.T) {
	num := int64(1234567890)
	b := util.Int64ToBytes(num)
	fmt.Printf("int64 2 []byte: %x\n", b)

	numBack := util.BytesToInt64(b)
	fmt.Printf("[]byte back to int64: %d\n", numBack)
}

func TestSnowflakeSeqGenerator_GenerateId(t *testing.T) {
	var dataCenterId, workId int64 = 1, 1
	generator, err := NewSnowflakeSeqGenerator(dataCenterId, workId)
	if err != nil {
		t.Error(err)
		return
	}
	var x, y string
	for i := 0; i < 100; i++ {
		y = generator.GenerateId("", "")
		if x == y {
			t.Errorf("x(%s) & y(%s) are the same", x, y)
		}
		x = y
	}
}

func TestRemovePunctuation(t *testing.T) {
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	mp := make(map[int64]struct{})
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println("Error initializing snowflake node:", err)
		return
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < 1000; j++ {
				//uuid := generator.GenerateId("", "")
				uuid := node.Generate()
				mu.Lock()
				if j == 0 {
					fmt.Printf("generated uid:%d\n", uuid)
				}
				if _, ok := mp[uuid.Int64()]; ok {
					fmt.Printf("error, id already exists....uuid is:%s\n", uuid)
				}
				mp[uuid.Int64()] = struct{}{}
				mu.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Printf("The len of map:%d\n", len(mp))
}
