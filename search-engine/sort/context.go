package sort

import (
	"Search-Engine/search-engine/model"
	"Search-Engine/search-engine/util"
	"fmt"
	"math"
	"sync"
)

// 用户进行搜索的 context
type SearchContext struct {
	mutex           sync.Mutex            // 多个 goroutine 在进行候选集添加的时候保护
	Query           string                // query语句，后续用于同标题计算最大公共子序列的长度
	CandidateDocIds []int64               `json:"candidateDocIds"` // 待选docId的数据集合
	CandidateItems  []model.CandidateItem // 候选集合的Item(Id & score)
}

func (c *SearchContext) DebugContext() {
	fmt.Printf("query:%s, len of candidate:%d\n, example candidate,docId :%d, score:%f, title:%s\n", c.Query,
		len(c.CandidateItems), c.CandidateItems[0].Id, c.CandidateItems[0].Score, c.CandidateItems[0].Title)
}

func (s *SearchContext) AddCandidate(docIdList *[]int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.CandidateDocIds = append(s.CandidateDocIds, *docIdList...)
}

func (s *SearchContext) PreProcess() {
	// 对于排序阶段的预先处理，例如计算 doc 中的词频出现量
	// 计算 docIdList 中每个id的重复次数，重复的doc越多说明对应doc中命中的term数量就越多
	hashSet := make(map[int64]struct{})
	for _, docId := range s.CandidateDocIds {
		itemsCnt := len(s.CandidateItems)
		_, exist := hashSet[docId]
		if itemsCnt == 0 || !exist {
			s.CandidateItems = append(s.CandidateItems, model.CandidateItem{
				Id:        docId,
				Score:     0,
				Frequency: 1,
			})
			hashSet[docId] = struct{}{}
		} else {
			s.CandidateItems[itemsCnt-1].Frequency++
		}
	}
}

func (s *SearchContext) AssignScores() {
	/*
		strategy:
			1. 针对于上面计算的待选集合的命中的频率, 命中频率比较高的分数肯定是更多的
			2. 直接判断一下标题和query的 最长公共子序列 的长度，较长的分数也要提上去
			3. (todo)可以计算多个 List 的交集，交集中的数据说明 全部命中，对应的分数也要提升上去
	*/
	// strategy2: 按照最长公共子序列的匹配程度来加分
	for index, item := range s.CandidateItems {
		candidateTitle, query := item.Title, s.Query
		lcsLen := util.CalculateLCS(candidateTitle, query)
		s.CandidateItems[index].Score += float64(lcsLen)
		s.CandidateItems[index].Score += float64(item.Frequency)
	}
}

// 对于待选结果集合进行去重，
func (s *SearchContext) DeduplicateResults() {
	candidate2bdeleted := make(map[int]struct{}) // 将重复的数据的下标放进来
	for i, v := range s.CandidateItems {
		// 这是已经被判断为重复的界面
		if _, exist := candidate2bdeleted[i]; exist {
			continue
		}
		lenx := len(v.Title)
		for j := i + 1; j < len(s.CandidateItems); j++ {
			if _, exist := candidate2bdeleted[j]; exist {
				continue
			}
			leny := len(s.CandidateItems[j].Title)
			abs := math.Abs(float64(lenx - leny))
			maxLen := math.Max(float64(lenx), float64(leny))
			// 长度差距过大，直接放过
			if abs/maxLen > 0.3 {
				continue
			}
			if float64(util.CalculateLCS(v.Title, s.CandidateItems[j].Title))/maxLen > 0.8 {
				candidate2bdeleted[j] = struct{}{}
			}
		}
	}
	var newCandidateItems []model.CandidateItem
	for i, v := range s.CandidateItems {
		if _, exist := candidate2bdeleted[i]; exist {
			continue
		}
		newCandidateItems = append(newCandidateItems, v)
	}
	fmt.Printf("进行了最终候选集中相似内容的去重，before:%d, after:%d\n", len(s.CandidateItems), len(newCandidateItems))
	s.CandidateItems = newCandidateItems
}
