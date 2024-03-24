package model

type SearchRequest struct {
	Query string `json:"query"` // 搜索时候键入的关键词
	Limit int    `json:"limit"` // 限制关键词的数量
}

// 键入搜索查询词之后的返回结果
type SearchResponse struct {
	TimeCost  float64             `json:"timeCost"` // 查询的耗时
	Terms     []string            `json:"terms"`    // 查询的关键词，主要是用于验证下正确性
	Documents []SearchResponseDoc `json:"documents"`
}

type SimpleSearchResponse struct {
	Terms      []string        `json:"terms"`
	Query      string          `json:"query"` // 将查询词也传递过去
	Candidates []CandidateItem `json:"candidates"`
}
