package container

import (
	"Search-Engine/search-engine/engine"
	tokenizer2 "Search-Engine/search-engine/words/tokenizer"
)

var GlobalContainer *Container

type Container struct {
	IndexStorageDir string // 索引数据存放的路径
	Tokenizer       *tokenizer2.Tokenizer
	engines         *engine.Engine
	ShardNum        int
	BufferNum       int
}

func InitGlobalContainer(tokenizer *tokenizer2.Tokenizer) {
	GlobalContainer = &Container{
		Tokenizer: tokenizer,
	}
}

func (c *Container) NewEngine() *engine.Engine {
	engine := &engine.Engine{
		IndexPath:             c.IndexStorageDir,
		Tokenizer:             c.Tokenizer,
		ShardNum:              c.ShardNum,
		BufferNum:             c.BufferNum,
		InvertIndexName:       "inverted_index",
		PositiveIndexName:     "positive_index",
		RepositoryStorageName: "repository_storage",
		TimeOut:               30,
	}
	engine.Init()
	return engine
}

func (c *Container) GetEngine() *engine.Engine {
	var engine *engine.Engine
	if c.engines == nil {
		engine = c.NewEngine()
	} else {
		engine = c.engines
	}
	return engine
}
