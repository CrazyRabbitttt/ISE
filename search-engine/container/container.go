package container

import (
	"Search-Engine/config"
	"Search-Engine/search-engine/engine"
	"Search-Engine/search-engine/reminder"
	tokenizer2 "Search-Engine/search-engine/words/tokenizer"
	"fmt"
	"os"
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
		Tokenizer:       tokenizer,
		IndexStorageDir: config.GlobalConfig.DB["default"].IndexStorageDir,
	}
}

func (c *Container) NewEngine() *engine.Engine {
	dbConfig := config.GlobalConfig.DB["default"]
	workDir, _ := os.Getwd()
	engine := &engine.Engine{
		IndexPath:             workDir + "/" + c.IndexStorageDir,
		Tokenizer:             c.Tokenizer,
		ShardNum:              c.ShardNum,
		BufferNum:             c.BufferNum,
		InvertIndexName:       dbConfig.InvertIndexName,
		PositiveIndexName:     dbConfig.PositiveIndexName,
		RepositoryStorageName: dbConfig.RepositoryStorageName,
		TimeOut:               dbConfig.TimeOut,
		TrieReminder:          reminder.NewTrie(),
	}
	fmt.Println("The index path is:", engine.IndexPath)
	engine.Init()
	return engine
}

func (c *Container) GetEngine() *engine.Engine {
	var engine *engine.Engine
	if c.engines == nil {
		fmt.Println("engine is nil, create")
		engine = c.NewEngine()
		c.engines = engine
		return c.engines
	} else {
		engine = c.engines
	}
	return engine
}
