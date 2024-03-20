package model

// 候选集合中的Item
type CandidateItem struct {
	Id        uint32  `json:"id"`        // docId
	Score     float64 `json:"score"`     // doc的分数、这个分数跟Query相关
	Frequency int     `json:"frequency"` // doc命中的 term 数量
	Title     string  `json:"title"`     // 特征中的 title 字段（或者说是叫做text）
	URL       string  `json:"url"`       // 特征中的URL，用于最终的结果展示
}
