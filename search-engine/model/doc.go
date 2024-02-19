package model

// 用于接收 Http Request请求的结构体
type IndexDoc struct {
	Key   uint32                 `json:"key"`   // 文档的唯一key
	Text  string                 `json:"terms"` // 分词后即是倒排索引的 terms
	Attrs map[string]interface{} `json:"attrs"` // 文档对应的特征(属性)
}

type RepositoryIndexDoc struct {
	*IndexDoc
	Terms []string `json:"terms"`
}

type SearchResponseDoc struct {
	IndexDoc
	Docs  []*RepositoryIndexDoc
	Score int
}
