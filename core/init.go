package core

import (
	tokenizer2 "Search-Engine/search-engine/words/tokenizer"
	"Search-Engine/web/router"
	"fmt"
)

// 进行初始化
func Initialize() {
	// 初始化分词器
	tokenizer := tokenizer2.NewTokenizer()
	str := "上海市的海上面有上海滩"
	words := tokenizer.Cut(str)
	for _, v := range words {
		fmt.Print(v, ",")
	}
	// 初始化路由
	r := router.InitRouter()
	r.Run()
}
