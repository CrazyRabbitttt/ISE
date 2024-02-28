package DataStructure

import (
	"sync"
	"time"
)

type MyConcurrentMap struct {
	sync.Mutex
	hashMap   map[int]int
	valueChan chan int
}

func (m *MyConcurrentMap) Put(k, v int) {

}

// 查询的时候 如果Key存在就返回，不存在就阻塞（直到存在或者说时间到了）
func (m *MyConcurrentMap) Get(k int, maxWaitingDuration time.Duration) (int, error) {
	v, exist := m.hashMap[k]
	if exist {
		return v, nil
	} else {
		select {
		case <-m.valueChan:
		}
	}
	return 0, nil
}
