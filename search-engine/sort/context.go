package sort

import (
	"Search-Engine/search-engine/model"
	"sync"
)

// 用户进行搜索的 context
type SearchContext struct {
	mutex           sync.Mutex            // 多个 goroutine 在进行候选集添加的时候保护
	CandidateDocIds []uint32              `json:"candidateDocIds"` // 待选docId的数据集合
	CandidateItems  []model.CandidateItem // 候选集合的Item(Id & score)
}

func (s *SearchContext) AddCandidate(docIdList *[]uint32) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.CandidateDocIds = append(s.CandidateDocIds, *docIdList...)
}

func (s *SearchContext) PreProcess() {
	// 计算 docIdList 中每个id的重复次数，重复的doc越多说明对应doc中命中的term数量就越多
	for _, docId := range s.CandidateDocIds {
		itemsCnt := len(s.CandidateItems)
		if itemsCnt == 0 || s.CandidateItems[itemsCnt-1].Id != docId { // 新增
			s.CandidateItems = append(s.CandidateItems, model.CandidateItem{
				Id:        docId,
				Score:     0,
				Frequency: 1,
			})
		} else {
			s.CandidateItems[itemsCnt-1].Frequency++
		}
	}
}

func (s *SearchContext) AssignScores() {
	/*
		strategy:
			1. 针对于上面计算的待选集合的命中的频率, 命中频率比较高的分数肯定是更多的
			2. 可以计算多个 List 的交集，交集中的数据说明 全部命中，对应的分数也要提升上去
			3. 直接判断一下标题和query的 最长公共子序列 的长度，较长的分数也要提上去
	*/

}
