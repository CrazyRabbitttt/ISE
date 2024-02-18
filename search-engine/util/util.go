package util

import "regexp"

func RemovePunctuation(word string) string {
	reg := regexp.MustCompile(`[[:punct:]]`) // 定义要删除的标点符号范围为所有标点符号
	return reg.ReplaceAllString(word, "")    // 将匹配到的标点符号全部替换成空字符串
}

func RemoveSpace(word string) string {
	reg := regexp.MustCompile(`\s+`)
	return reg.ReplaceAllString(word, "")
}
