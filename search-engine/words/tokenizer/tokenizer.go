package tokenizer

import (
	"Search-Engine/search-engine/util"
	"github.com/wangbin/jiebago"
	"strings"
)

// 用于进行 Word 的分词处理

type Tokenizer struct {
	seg jiebago.Segmenter
}

func NewTokenizer() *Tokenizer {
	//workdir, err := os.Getwd()
	//if err != nil {
	//	panic(err)
	//}
	tokenizer := &Tokenizer{}
	//err := tokenizer.seg.LoadDictionary("../../../data/dictionary.txt")
	err := tokenizer.seg.LoadDictionary("data/dictionary.txt")
	if err != nil {
		panic(err)
	}
	return tokenizer
}

func (c *Tokenizer) Cut(word string) []string {
	var resWords []string
	// 忽略大小写
	word = strings.ToLower(word)
	// 去掉标点符号
	word = util.RemovePunctuation(word)
	// 去掉空白符号
	word = util.RemoveSpace(word)
	// wordSet
	wordSet := make(map[string]struct{})
	receiveChannel := c.seg.CutForSearch(word, true)
	for str := range receiveChannel {
		_, exist := wordSet[str]
		if !exist {
			wordSet[str] = struct{}{}
			resWords = append(resWords, str)
		}
	}
	return resWords
}
