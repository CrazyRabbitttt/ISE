package core

import (
	"Search-Engine/config"
	"Search-Engine/search-engine/container"
	tokenizer2 "Search-Engine/search-engine/words/tokenizer"
	"Search-Engine/web/router"
	"Search-Engine/web/service"
)

// 进行初始化
func Initialize() {
	// 初始化全局的配置文件
	config.InitConfig()
	// 初始化分词器
	tokenizer := tokenizer2.NewTokenizer()
	// 初始化全局的 Container
	container.InitGlobalContainer(tokenizer)
	// 初始化业务逻辑
	service.InitService()
	// 初始化路由
	r := router.InitRouter()
	// 将一些测试用的query词加到 Trie 树中去

	r.Run()
}
