package engine

import (
	"Search-Engine/search-engine/model"
	"Search-Engine/search-engine/storage"
	tokenizer2 "Search-Engine/search-engine/words/tokenizer"
	"sync"
)

type Engine struct {
	IndexPath            string                    // 索引数据存放的路径
	Tokenizer            *tokenizer2.Tokenizer     // 分词器
	InvertedIndexStorage *[]storage.LeveldbStorage // 倒排索引 term ==> id list
	PositiveIndexStorage *[]storage.LeveldbStorage // 正排索引 id ==> terms
	RepositoryStorage    *[]storage.LeveldbStorage // id ==> terms + attrs
	ShardNum             int                       // shard数目（默认10）
	BufferNum            int                       // 每个shard存放的kv对的数量（默认1000）
	AddIndexDocChan      []chan *model.IndexDoc    // 用于接收需要构建的 IndexDoc 的 channel

	InvertIndexName       string // 三个索引文件的前缀名称
	PositiveIndexName     string
	RepositoryStorageName string

	TimeOut int64 // 超时时间

	wg sync.WaitGroup // 用于同步 初始化 和 执行索引操作的 goroutine
}

func (engine *Engine) Init() {
	engine.wg.Add(1)
	defer engine.wg.Done()

	if engine.ShardNum == 0 {
		engine.ShardNum = 10
	}
	if engine.BufferNum == 0 {
		engine.BufferNum = 1000
	}
	//engine.AddIndexDocChan = make([]chan *model.IndexDoc, engine.ShardNum)
	//for i := 0; i < engine.ShardNum; i++ {
	//	engine.AddIndexDocChan[i] = make(chan *model.IndexDoc, engine.BufferNum)
	//	go engine.AddIndexDocLoop(engine.AddIndexDocChan[i])
	//}

	// 初始化三个 Leveldb 的访问
	//s, err := storage.NewStorage()

}

func (engine *Engine) AddIndexDocLoop(worker chan *model.IndexDoc) {
	for {
		indexDoc := <-worker // 试图从channel中读取待添加的Doc，没有的话就阻塞在这里
		engine.AddIndexDoc(indexDoc)
	}
}

func (e *Engine) AddIndexDoc(indexDoc *model.IndexDoc) {
	// 需要等待初始化完成才能添加索引
	e.wg.Wait()

	//tokenizer := e.Tokenizer
	//terms := tokenizer.Cut(indexDoc.Text)
	//docId := indexDoc.Key

}
