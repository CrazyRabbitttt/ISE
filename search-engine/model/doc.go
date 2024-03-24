package model

import (
	"bytes"
	"encoding/gob"
)

// 用于接收 Http Request请求的结构体
type IndexDoc struct {
	Key   int64             `json:"key"`   // 文档的id（这里不再需要业务方提供docId，由server端生成）
	Text  string            `json:"terms"` // 分词后即是倒排索引的 terms
	Attrs map[string]string `json:"attrs"` // 文档对应的特征(属性)
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

func (i *IndexDoc) Encode() ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	// 对结构体进行编码
	err := encoder.Encode(i)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (i *IndexDoc) Decode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)

	// 对字节切片进行解码
	err := decoder.Decode(i)
	if err != nil {
		return err
	}

	return nil
}
