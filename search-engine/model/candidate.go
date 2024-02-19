package model

// 候选集合中的Item
type CandidateItem struct {
	Id        uint32  `json:"id"`
	Score     float64 `json:"score"`
	Frequency int     `json:"frequency"` // doc命中的 term 数量
}
