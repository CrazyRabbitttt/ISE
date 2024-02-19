package engine

import (
	"Search-Engine/search-engine/model"
	"Search-Engine/search-engine/sort"
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

	mutex sync.Mutex
	wg    sync.WaitGroup // 用于同步 初始化 和 执行索引操作的 goroutine
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
		// 倒排索引
		s, err := storage.NewStorage(engine.IndexPath+"/"+invertIndexFileName, engine.TimeOut)
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
		fmt.Println("成功添加IndexDoc到Channel中, docId:", indexDoc.Key, "text", indexDoc.Text, "attrs", indexDoc.Attrs)
		engine.AddIndexDoc2Engine(indexDoc)
	}
}

func (e *Engine) AddIndexDoc2Chan(indexDoc *model.IndexDoc) {
	// 需要等待初始化完成才能添加索引
	docId := indexDoc.Key
	e.DocumentCnt++
	e.AddIndexDocChan[e.GetShardNumByDocId(docId)] <- indexDoc
}

func (e *Engine) AddIndexDoc2Engine(indexDoc *model.IndexDoc) {
	// 等待初始化完成后进行数据添加
	e.wg.Wait()

	docId := indexDoc.Key
	terms := e.Tokenizer.Cut(indexDoc.Text)
	/*
		  倒排索引：
				如果说 docId 之前是不存在的，那么就是纯增加。
				如果 doc 之前是存在过的，那么就可能涉及到更改
				old:
					docId:123
					terms: 苹果、香蕉、梨子、火龙果
				new:
					docId:123
					terms: 苹果、葡萄、栗子、火龙果

				   针对于 新增的 [葡萄、栗子] 对应的倒排链中，就应该加上 docID
			       针对于 移除的 [香蕉、梨子] 对应的倒排链中，就应该去掉 docID
		  正排索引：
				针对于正排索引 (docId ===> [terms], docId ===> [terms + document]) 来说，
				其实直接在 Leveldb 中 Set 就可以了, 因为即使docId曾经作为key存在于db中，leveldb
				也会将value给替换掉。
	*/
	terms2bRemoved, terms2bInserted := e.PrepareForHandle(terms, docId) // 内置了对于DB的handle， 需要进行加锁🔒
	fmt.Printf("The len of remove:%d, the len of insert:%d", len(terms2bRemoved), len(terms2bInserted))
	// 倒排索引：删除索引
	for _, value := range terms2bRemoved {
		e.RemoveDocIdInInvertIndex(value, docId)
	}
	// 倒排索引：新增索引
	for _, value := range terms2bInserted {
		e.AddDocIdInInvertIndex(value, docId)
	}

	// 更新正排索引
	e.AddIndexDoc2PositiveIndex(indexDoc, terms)
}

func (e *Engine) AddIndexDoc2PositiveIndex(indexDoc *model.IndexDoc, terms []string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	index := e.GetShardNumByDocId(indexDoc.Key)
	positiveIndex := e.PositiveIndexStorage[index]
	reposIndex := e.RepositoryStorage[index]

	repos := &model.RepositoryIndexDoc{
		IndexDoc: indexDoc,
		Terms:    terms,
	}

	// id ===> [terms]
	positiveIndex.Set(util.Uint32ToBytes(indexDoc.Key), util.Encoder(terms))
	// id ===> [terms + attrs]
	reposIndex.Set(util.Uint32ToBytes(indexDoc.Key), util.Encoder(repos))
}

func (e *Engine) PrepareForHandle(terms []string, docId uint32) ([]string, []string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	docIndex := e.GetShardNumByDocId(docId)
	positiveIndex := e.PositiveIndexStorage[docIndex] // 从里面取到的list是docId ===> term list
	buf, exist := positiveIndex.Get(util.Uint32ToBytes(docId))
	var terms2bRemoved []string
	var terms2bInserted []string
	if !exist { // docId 本身就是不存在的，那么直接添加索引数据
		terms2bInserted = terms
	} else { // docId 本身是存在的，那么本次传递过来的数据可能涉及到倒排索引的更新(新建、删除)
		var oldTermList []string
		util.Decoder(buf, &oldTermList)
		// 需要被删除掉的 terms
		for _, oldTerm := range oldTermList {
			_, exist := util.ExistInArrayString(terms, oldTerm)
			if !exist {
				terms2bRemoved = append(terms2bRemoved, oldTerm)
			}
		}
		// 需要新增的 terms
		for _, newTerm := range terms {
			_, exist := util.ExistInArrayString(oldTermList, newTerm)
			if !exist {
				terms2bInserted = append(terms2bInserted, newTerm)
			}
		}
	}
	return terms2bRemoved, terms2bInserted
}

// 给到 term 对应的倒排链中添加一个 docId
func (e *Engine) AddDocIdInInvertIndex(term string, docId uint32) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	fmt.Println("执行倒排索引结构的增添, term:", term, ", docId:", docId)
	var docIdList = make([]uint32, 0)
	invertIndex := e.InvertedIndexStorage[e.GetShardNumByTerm(term)]
	buf, exist := invertIndex.Get([]byte(term))
	if !exist {
		docIdList = append(docIdList, docId)
		fmt.Println("AddInvertIndex function, 没有%s构建的倒排索引", term)
	} else {
		util.Decoder(buf, &docIdList)
		if _, exist := util.ExistInArrayUint32(docIdList, docId); !exist {
			docIdList = append(docIdList, docId)
		}
		fmt.Println("AddInvertIndex function, 将%d添加到%s对应的倒排拉链中", docId, term)
	}
	// 将更新后的 docIdList 设置到 db 中
	invertIndex.Set([]byte(term), util.Encoder(docIdList))
}

func (e *Engine) RemoveDocIdInInvertIndex(term string, docId uint32) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	var docIdList = make([]uint32, 0)
	invertIndex := e.InvertedIndexStorage[e.GetShardNumByTerm(term)]
	buf, exist := invertIndex.Get([]byte(term))
	if exist {
		util.Decoder(buf, &docIdList)
		// 将 id 从 list 中删除
		docIdList = util.RemoveUint32ValueInArray(docIdList, docId)
		if len(docIdList) == 0 {
			// 这个倒排索引已经空了，直接删掉
			if err := invertIndex.Delete([]byte(term)); err != nil {
				panic(err)
			}
		} else {
			invertIndex.Set([]byte(term), util.Encoder(docIdList))
		}
	}
}

func (e *Engine) Search(request *model.SearchRequest) (*model.SearchResponse, error) {
	searchContext := &sort.SearchContext{} // 本次搜索的上下文（包括待选数据集等）
	// 1. 首先对于 query 进行分词处理
	terms := e.Tokenizer.Cut(request.Query)
	termCnt := len(terms)
	// 2. 查询倒排索引，获取 terms 对应的 docIdList 作为候选结果集，后续排序啥的用
	wg := &sync.WaitGroup{}
	for i := 0; i < termCnt; i++ {
		go e.AddDocIdList2ContextByTerm(terms[i], searchContext, wg)
		wg.Done()
	}
	wg.Wait() // 等待多线程完成对于不同分片的 候选集 的添加

	// 3. Preprocessing 预处理数据, 获得待选集doc命中的term数量
	searchContext.PreProcess()
	// 4. AssignScore 赋分数
	searchContext.AssignScores()
}

func (e *Engine) AddDocIdList2ContextByTerm(term string, context *sort.SearchContext, wg *sync.WaitGroup) {
	defer wg.Done()
	index := e.GetShardNumByTerm(term)
	invertIndex := e.InvertedIndexStorage[index]
	var docIdList []uint32
	buf, exist := invertIndex.Get([]byte(term))
	if exist {
		util.Decoder(buf, &docIdList)
		context.AddCandidate(&docIdList)
	}
}

func (e *Engine) GetShardNumByDocId(docId uint32) int {
	return int(docId % uint32(e.ShardNum))
}

func (e *Engine) GetShardNumByTerm(term string) int {
	// murmur hash ==> hash code
	hashCode := util.String2Int(term)
	return int(hashCode % uint32(e.ShardNum))
}
