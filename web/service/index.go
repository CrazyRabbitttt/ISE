package service

import (
	"Search-Engine/search-engine/container"
	"Search-Engine/search-engine/model"
	"Search-Engine/search-engine/util"
	"fmt"
)

type IndexService struct {
	container *container.Container
}

func NewIndexService() *IndexService {
	return &IndexService{
		container: container.GlobalContainer,
	}
}

func (s *IndexService) DebugPositiveIndex(doc *model.IndexDoc) {
	docId := doc.Key
	engine := s.container.GetEngine()

	terms := engine.Tokenizer.Cut(doc.Text)
	shardNum := engine.GetShardNumByDocId(docId)
	repo := engine.RepositoryStorage[shardNum]
	repos := &model.RepositoryIndexDoc{
		IndexDoc: doc,
		Terms:    terms,
	}
	buf, _ := util.Encoder(repos)
	repo.Set(util.Int64ToBytes(docId), buf)

	res := new(model.RepositoryIndexDoc)
	b, _ := repo.Get(util.Int64ToBytes(docId))
	util.Decoder(b, &res)

	fmt.Printf("After decode: docId:%d, text:%s, title:%s\n", res.Key, res.Text, res.Attrs["title"])

}

func (s *IndexService) DebugIndex(doc *model.IndexDoc) error {
	docId := doc.Key
	engine := s.container.GetEngine()

	terms := engine.Tokenizer.Cut(doc.Text)

	var docIdList = make([]int64, 0)
	docIdList = append(docIdList, docId)
	for i, v := range terms {
		shardNum := engine.GetShardNumByTerm(v)
		invertStorage := engine.InvertedIndexStorage[shardNum]
		fmt.Printf("tik tok:%d, term:%s, shard number:%d\n", i, v, shardNum)
		if _, exist := invertStorage.Get([]byte(v)); !exist {
			b, _ := util.Encoder(docIdList)
			invertStorage.Set([]byte(v), b)
		}
	}

	for i, v := range terms {
		shardNum := engine.GetShardNumByTerm(v)
		invertStorage := engine.InvertedIndexStorage[shardNum]
		fmt.Printf("tok tik:%d, term:%s, shard number:%d\n", i, v, shardNum)
		if b, exist := invertStorage.Get([]byte(v)); exist {
			var l = make([]int64, 0)
			util.Decoder(b, &l)
			fmt.Println("doc id list:", l)
		}
	}
	return nil
}

func (s *IndexService) AddIndexDoc(doc *model.IndexDoc) error {
	// 往Engine对应的chan中添加数据
	go s.container.GetEngine().AddIndexDoc2Chan(doc)
	return nil
}
