package engine

import (
	"Search-Engine/search-engine/model"
	"Search-Engine/search-engine/reminder"
	"Search-Engine/search-engine/sort"
	"Search-Engine/search-engine/storage"
	"Search-Engine/search-engine/util"
	tokenizer2 "Search-Engine/search-engine/words/tokenizer"
	"fmt"
	stdsort "sort"
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
	TrieReminder          *reminder.Trie

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
	fmt.Printf("添加索引，到倒排索引中需要删除的 terms 的长度:%d, 倒排索引中需要添加的 terms 的长度:%d\n", len(terms2bRemoved), len(terms2bInserted))
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

	docId := indexDoc.Key
	index := e.GetShardNumByDocId(docId)
	fmt.Printf("更新正排索引, docId:%d, 对应的索引下标:%d, ", docId, index)
	positiveIndex := e.PositiveIndexStorage[index]
	reposIndex := e.RepositoryStorage[index]

	repos := &model.RepositoryIndexDoc{
		IndexDoc: indexDoc,
		Terms:    terms,
	}
	fmt.Printf("Before handle it in index, data: %T\n", terms)

	// id ===> [terms]
	buf1, err := util.Encoder(terms)
	if err != nil {
		fmt.Printf("when encode terms, occur Error:%v\n", err)
	} else {
		fmt.Printf("执行了 positive encode, 但是正常\n")
	}
	positiveIndex.Set(util.Int64ToBytes(docId), buf1)
	fmt.Printf("Type of repo:%T,Value of repo, key:%d, temrs[0]:%s\n", repos, repos.Key, repos.Terms[0])
	// id ===> [terms + attrs]
	encodedData, err := indexDoc.Encode()
	if err != nil {
		fmt.Printf("when encode repo,  occur Error:%v\n", err)
		fmt.Printf("docId:%d, text:%s, terms[0]:%s, title:%s\n", repos.Key, repos.Text, repos.Terms[0], repos.Attrs["title"])
	} else {
		fmt.Printf("执行了 repo encode, 但是正常\n")
	}
	//buf2, err := util.Encoder(repos)
	//if err != nil {
	//	fmt.Printf("when encode repo,  occur Error:%v\n", err)
	//	fmt.Printf("docId:%d, text:%s, terms[0]:%s, title:%s\n", repos.Key, repos.Text, repos.Terms[0], repos.Attrs["title"])
	//} else {
	//	fmt.Printf("执行了 repo encode, 但是正常")
	//}
	reposIndex.Set(util.Int64ToBytes(docId), encodedData)
}

func (e *Engine) PrepareForHandle(terms []string, docId int64) ([]string, []string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	docIndex := e.GetShardNumByDocId(docId)
	positiveIndex := e.PositiveIndexStorage[docIndex] // 从里面取到的list是docId ===> term list
	buf, exist := positiveIndex.Get(util.Int64ToBytes(docId))
	var terms2bRemoved []string
	var terms2bInserted []string
	if !exist { // docId 本身就是不存在的，那么直接添加索引数据
		terms2bInserted = terms
	} else { // docId 本身是存在的，那么本次传递过来的数据可能涉及到倒排索引的更新(新建、删除)
		oldTermList := make([]string, 0)
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
func (e *Engine) AddDocIdInInvertIndex(term string, docId int64) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	fmt.Println("执行倒排索引结构的增添, term:", term, ", docId:", docId)
	var docIdList = make([]int64, 0)
	invertIndex := e.InvertedIndexStorage[e.GetShardNumByTerm(term)]
	buf, exist := invertIndex.Get([]byte(term))
	if !exist {
		docIdList = append(docIdList, docId)
		fmt.Printf("AddInvertIndex function, 没有%s构建的倒排索引\n", term)
	} else {
		util.Decoder(buf, &docIdList)
		if _, exist := util.ExistInArrayUint32(docIdList, docId); !exist {
			docIdList = append(docIdList, docId)
		}
		fmt.Printf("AddInvertIndex function, 将%d添加到%s对应的倒排拉链中\n", docId, term)
	}
	// 将更新后的 docIdList 设置到 db 中
	buf1, err := util.Encoder(docIdList)
	if err != nil {
		fmt.Printf("when encode docIdList, error occur:%v\n", err)
	}
	invertIndex.Set([]byte(term), buf1)
}

func (e *Engine) RemoveDocIdInInvertIndex(term string, docId int64) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	var docIdList = make([]int64, 0)
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
			buf1, err := util.Encoder(docIdList)
			if err != nil {
				fmt.Printf("when encode docIdlist, error occr:%v\n", err)
			}
			invertIndex.Set([]byte(term), buf1)
		}
	}
}

func (e *Engine) Search(request *model.SearchRequest) (*model.SimpleSearchResponse, error) {
	searchContext := &sort.SearchContext{ // 本次搜索的上下文（包括待选数据集等）
		Query: request.Query,
	}
	// 1. 首先对于 query 进行分词处理
	terms := e.Tokenizer.Cut(request.Query)
	termCnt := len(terms)
	// 2. 查询倒排索引，获取 terms 对应的 docIdList 作为候选结果集，后续排序啥的用
	wg := &sync.WaitGroup{}
	wg.Add(termCnt)
	for i := 0; i < termCnt; i++ {
		go e.AddDocIdList2ContextByTerm(terms[i], searchContext, wg) // 类似于查询倒排索引这一步
	}
	wg.Wait() // 等待多线程完成对于不同分片的 候选集 的添加
	fmt.Println(searchContext.CandidateDocIds)
	// 3. Preprocessing 预处理数据, 获得待选集doc命中的term数量
	searchContext.PreProcess()
	// 4. 拿到 docId 对应的一些特征，这里其实类似于查询 正排索引 这一步, 将候选集的特征进行 Enrich
	e.AddAttrs2ContextByDocId(searchContext, wg)
	// 5. AssignScore 赋分数
	searchContext.AssignScores()
	// 6. 排序
	stdsort.Sort(model.CandidateItemSlice(searchContext.CandidateItems))
	// 7.再次进行截断
	if len(searchContext.CandidateItems) > 25 {
		fmt.Printf("最终排序后进行截断, before:%d, after:25\n", len(searchContext.CandidateItems))
		searchContext.CandidateItems = searchContext.CandidateItems[:25]
	}
	// 6. 这里对于 query 词加到 Trie 树中，用于关键词提示(todo : 后面需要在启动的时候将一些热搜词直接初始化好)
	e.TrieReminder.Add(request.Query)
	response := &model.SimpleSearchResponse{
		Terms:      terms,
		Candidates: searchContext.CandidateItems,
	}
	return response, nil
}

func (e *Engine) AddDocIdList2ContextByTerm(term string, context *sort.SearchContext, wg *sync.WaitGroup) {
	defer wg.Done()
	index := e.GetShardNumByTerm(term)
	invertIndex := e.InvertedIndexStorage[index]
	var docIdList []int64
	buf, exist := invertIndex.Get([]byte(term))
	if exist {
		util.Decoder(buf, &docIdList)
		context.AddCandidate(&docIdList)
	}
}

func (e *Engine) AddAttrs2ContextByDocId(context *sort.SearchContext, wg *sync.WaitGroup) {
	// 获得 docId 对应的文档库，拿到一些特征（例如Title、URL、作者、文档的描述等等）
	// 这里可以开启多个 goroutine 同时获取 doc 对应的特征
	var newCandidateItems []model.CandidateItem
	var newCandidateIds []int64
	wg.Add(len(context.CandidateItems))
	for i, item := range context.CandidateItems {
		docId := item.Id
		go e.GetAttrsFromStorageByDocId(docId, &context.CandidateItems[i], wg)
	}
	wg.Wait()
	for i, v := range context.CandidateItems {
		if v.Title == "" || v.URL == "" || v.URL == " " {
			continue
		}
		newCandidateItems = append(newCandidateItems, v)
		newCandidateIds = append(newCandidateIds, context.CandidateDocIds[i])
	}
	if len(newCandidateItems) != len(newCandidateIds) {
		fmt.Printf("异常！！！！！经过过滤的结果集的docList和Candidate长度不一样！！！！！%d:%d\n",
			len(newCandidateIds), len(newCandidateItems))
	}
	if len(newCandidateItems) != len(context.CandidateItems) {
		fmt.Printf("对于候选集合中的空结果集进行了筛选，origin:%d, after:%d\n", len(context.CandidateItems), len(newCandidateItems))
		context.CandidateDocIds = newCandidateIds
		context.CandidateItems = newCandidateItems
	}
	// 算是首次召回的截断
	if len(context.CandidateItems) > 60 {
		context.CandidateItems = context.CandidateItems[:60]
	}
}

// GetAttrsFromStorageByDocId 函数获取 doc 对应的文档，从中取出来特征并且赋值给传入的 候选集 的字段中
func (e *Engine) GetAttrsFromStorageByDocId(docId int64, candidateItem *model.CandidateItem, wg *sync.WaitGroup) {
	e.mutex.Lock()
	e.mutex.Unlock()
	defer wg.Done()

	shardIndex := e.GetShardNumByDocId(docId)
	storageIndex := e.RepositoryStorage[shardIndex]

	buf, exist := storageIndex.Get(util.Int64ToBytes(docId))
	if exist {
		fmt.Printf("存在docId对应的attr, docId:%d\n", docId)
		//repos := new(model.RepositoryIndexDoc)
		var repos model.IndexDoc
		//util.Decoder(buf)
		err := repos.Decode(buf)
		if err != nil {
			fmt.Println("Error decoding:", err)
			return
		}
		attrs := repos.Attrs
		candidateItem.Title = attrs["title"]
		candidateItem.URL = attrs["page_url"]
		candidateItem.Description = attrs["description"]
		candidateItem.KeyWords = attrs["keywords"]
		fmt.Printf("Assign url, docId:%d, url:%s\n", candidateItem.Id, candidateItem.URL)
		//titleInterface := attrs["title"]
		//urlInterface := attrs["page_url"]
		//if pageUrl, ok := urlInterface.(string); ok {
		//	candidateItem.URL = pageUrl
		//	fmt.Printf("Assign url, docId:%d, url:%s\n", candidateItem.Id, candidateItem.URL)
		//} else {
		//	fmt.Printf("There is no url in attrs, docId:%d\n", candidateItem.Id)
		//}
		//if title, ok := titleInterface.(string); ok {
		//	candidateItem.Title = title
		//	fmt.Printf("Assign title, docId:%d, title:%s\n", candidateItem.Id, candidateItem.Title)
		//}
	} else {
		fmt.Printf("There is no doc in storage, doc id:%d\n", docId)
	}
}

func (e *Engine) GetShardNumByDocId(docId int64) int {
	return int(docId % int64(e.ShardNum))
}

func (e *Engine) GetShardNumByDocIdStr(docId string) int {
	hashCode := util.String2Int(docId)
	return int(hashCode % uint32(e.ShardNum))
}

func (e *Engine) GetShardNumByTerm(term string) int {
	// murmur hash ==> hash code
	hashCode := util.String2Int(term)
	return int(hashCode % uint32(e.ShardNum))
}

func (e *Engine) SearchRemind(query string) ([]string, error) {
	var res []string
	nodes := e.TrieReminder.PrefixSearch(query, 10)
	for _, node := range nodes {
		res = append(res, node.Data.(string))
	}
	return res, nil
}

func (e *Engine) InitReminder(querys []string) {
	for index, query := range querys {
		fmt.Printf("init number [%d] query to Trie Tree\n", index)
		e.TrieReminder.Add(query)
	}
}
