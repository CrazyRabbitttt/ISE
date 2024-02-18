package main

import (
	"Search-Engine/search-engine/words/tokenizer"
	"fmt"
)

func main() {
	fmt.Println("Hello World")
	tokenize := tokenizer.NewTokenizer()
	originStr := "This is my firstsear''chengineP roject, sfd,哈哈哈，有什么问'题么，你是什么动物呢    ..,,你是小学生吗"
	words := tokenize.Cut(originStr)
	for _, word := range words {
		fmt.Println(word)
	}
}
