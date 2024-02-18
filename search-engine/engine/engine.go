package engine

import (
	"Search-Engine/search-engine/model"
	"Search-Engine/search-engine/storage"
	"Search-Engine/search-engine/util"
	tokenizer2 "Search-Engine/search-engine/words/tokenizer"
	"fmt"
	"sync"
)

type Engine struct {
	IndexPath            string                    // 索引数据存放的路径
	Tokenizer            *tokenizer2.Tokenizer     // 分词器
	InvertedIndexStorage []*storage.LeveldbStorage // 倒排索引 term ==> id list
	PositiveIndexStorage []*storage.LeveldbStorage // 正排索引 id ==> terms
	RepositoryStorage    []*storage.LeveldbStorage // id ==> terms + attrs
	ShardNum             int                       // shard数目（默认10）
	BufferNum            int                       // 每个shard存放的kv对的数量（默认1000）
	AddIndexDocChan      []chan *model.IndexDoc    // 用于接收需要构建的 IndexDoc 的 channel

	InvertIndexName       string // 三个索引文件的前缀名称
	PositiveIndexName     string
	RepositoryStorageName string
	DocumentCnt           int
	TimeOut               int64 // 超时时间

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
	engine.AddIndexDocChan = make([]chan *model.IndexDoc, engine.ShardNum)
	for i := 0; i < engine.ShardNum; i++ {
		engine.AddIndexDocChan[i] = make(chan *model.IndexDoc, engine.BufferNum)
		go engine.AddIndexDocLoop(engine.AddIndexDocChan[i])
		// 初始化三个 Leveldb 的访问
		invertIndexFileName := fmt.Sprintf("%s_%d", engine.InvertIndexName, i)
		positiveIndexFileName := fmt.Sprintf("%s_%d", engine.PositiveIndexName, i)
		repositoryStorageName := fmt.Sprintf("%s_%d", engine.RepositoryStorageName, i)
		fmt.Println("The invertIndexfileName:", invertIndexFileName)
		// 倒排索引
		//s, err := storage.NewStorage(engine.IndexPath+"/"+invertIndexFileName, engine.TimeOut)
		s, err := storage.NewStorage("./data/index_data/"+invertIndexFileName, engine.TimeOut)
		if err != nil {
			panic(err)
		}
		engine.InvertedIndexStorage = append(engine.InvertedIndexStorage, s)
		// 正排索引
		sP, errP := storage.NewStorage(engine.IndexPath+"/"+positiveIndexFileName, engine.TimeOut)
		if errP != nil {
			panic(err)
		}
		engine.PositiveIndexStorage = append(engine.PositiveIndexStorage, sP)
		// 文档存储
		sS, errS := storage.NewStorage(engine.IndexPath+"/"+repositoryStorageName, engine.TimeOut)
		if errS != nil {
			panic(errS)
		}
		engine.RepositoryStorage = append(engine.RepositoryStorage, sS)
		fmt.Println("End of new engine function.")
	}

}

func (engine *Engine) AddIndexDocLoop(worker chan *model.IndexDoc) {
	for {
		indexDoc := <-worker // 试图从channel中读取待添加的Doc，没有的话就阻塞在这里
		engine.AddIndexDoc2Engine(indexDoc)
	}
}

func (e *Engine) AddIndexDoc2Chan(indexDoc *model.IndexDoc) {
	// 需要等待初始化完成才能添加索引
	docId := indexDoc.Key
	e.DocumentCnt++
	e.AddIndexDocChan[e.GetShardNum(docId)] <- indexDoc
}

func (e *Engine) AddIndexDoc2Engine(indexDoc *model.IndexDoc) {
	// 等待初始化完成后进行数据添加
	e.wg.Wait()

	docId := indexDoc.Key
	terms := e.Tokenizer.Cut(indexDoc.Text)
	/*
			如果说 docId 之前是不存在的，那么就是纯增加。
			如果 doc 之前是存在过的，那么就可能涉及到更改
			old:
				docId:123
				terms: 苹果、香蕉、梨子、火龙果

			new:
				docId:123
				terms: 苹果、葡萄、栗子、火龙果

			那么针对于 新增的 [葡萄、栗子] 对应的倒排链中，就应该加上 docID
		       针对于 移除的 [香蕉、梨子] 对应的倒排链中，就应该去掉 docID
	*/

	invertIndex := e.InvertedIndexStorage
	// todo：

}

func (e *Engine) GetShardNum(docId uint32) int {
	return int(docId % uint32(e.ShardNum))
}

func (e *Engine) GetShardNumByTerm(term string) int {
	// murmur hash ==> hash code
	hashCode := util.String2Int(term)
	return int(hashCode % uint32(e.ShardNum))
}
