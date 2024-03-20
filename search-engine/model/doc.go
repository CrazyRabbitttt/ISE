package model

// 用于接收 Http Request请求的结构体
type IndexDoc struct {
	Key   uint32                 `json:"key"`   // 文档的id（这里不再需要业务方提供docId，由server端生成）
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

type InitTrie struct {
	Querys map[string]interface{} `json:"querys"`
}
