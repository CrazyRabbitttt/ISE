package model

// 候选集合中的Item
type CandidateItem struct {
	Id          int64   `json:"id"`          // docId
	Score       float64 `json:"score"`       // doc的分数、这个分数跟Query相关
	Frequency   int     `json:"frequency"`   // doc命中的 term 数量
	Title       string  `json:"title"`       // 特征中的 title 字段（或者说是叫做text）
	URL         string  `json:"url"`         // 特征中的URL，用于最终的结果展示
	Description string  `json:"description"` // 对于整个文档的描述信息
	KeyWords    string  `json:"keywords"`    // 整个文档中的关键词
}

type CandidateItemSlice []CandidateItem

func (c CandidateItemSlice) Len() int {
	return len(c)
}

func (c CandidateItemSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c CandidateItemSlice) Less(i, j int) bool {
	return c[i].Score > c[j].Score
}
