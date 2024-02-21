package reminder

import (
	"container/heap"
	"fmt"
)

type TrieNodeHeap []*TrieNode

func (h TrieNodeHeap) Len() int {
	return len(h)
}

func (h TrieNodeHeap) Less(i, j int) bool {
	return h[i].passCnt > h[j].passCnt
}

func (h *TrieNodeHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *TrieNodeHeap) Push(node interface{}) {
	*h = append(*h, node.(*TrieNode))
}

func (h *TrieNodeHeap) Pop() interface{} {
	res := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return res
}

type TrieNode struct {
	childNode map[rune]*TrieNode

	char      rune        // 节点的字符
	isTermEnd bool        // 是否是单词的末尾
	Data      interface{} // 单词的结尾存放数据（可以是这个单词的string）
	passCnt   int32       // 经过这个节点的数量 (这个节点必须是term的结尾)
	queryCnt  int32       // 本单词被查询的次数
}

type Trie struct {
	root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{
		root: newTrieNode(' '),
	}
}

func newTrieNode(ch rune) *TrieNode {
	return &TrieNode{
		char:      ch,
		childNode: make(map[rune]*TrieNode, 10),
		isTermEnd: false,
		Data:      nil,
		passCnt:   0,
		queryCnt:  0,
	}
}

func (t *Trie) Add(str string) {
	var curNode = t.root
	charLists := []rune(str)
	for _, ch := range charLists {
		nextNode, exist := curNode.childNode[ch]
		if !exist {
			nextNode = newTrieNode(ch)
			curNode.childNode[ch] = nextNode
		}
		if nextNode.isTermEnd {
			nextNode.passCnt++
		}
		curNode = nextNode
	}
	curNode.isTermEnd = true
	curNode.passCnt++
	curNode.Data = str // 末尾节点记录字符串，便于直接返回

	fmt.Printf("add string[%s] to trie\n", str)
}

// prefixSearch: 给定前缀查询词, 获取提示词
func (t *Trie) PrefixSearch(query string, limit int) (smallHeap TrieNodeHeap) {
	curNode := t.root
	charLists := []rune(query)
	heap.Init(&smallHeap)
	for _, ch := range charLists {
		nextNode, exist := curNode.childNode[ch]
		if !exist { // 前缀并不存在于 Trie 树中，返回的是一个空的
			fmt.Printf("prefix [%s] not in trie tree\n", query)
			return
		}
		curNode = nextNode
	}
	//fmt.Println("debug1")
	// curNode 目前是query中最后一个字符对应的节点
	rangeQueue := []*TrieNode{curNode}
	for len(rangeQueue) > 0 { // while !(queue.empty())
		fmt.Println("len of queue:", len(rangeQueue))
		var tmpQueue []*TrieNode
		for _, node := range rangeQueue {
			// 将 Term 放进小顶堆中
			if node.isTermEnd == true && node != curNode {
				heap.Push(&smallHeap, node)
			}
			for _, v := range node.childNode {
				tmpQueue = append(tmpQueue, v)
			}
		}
		rangeQueue = tmpQueue
	}
	//fmt.Println("debug2")
	if smallHeap.Len() > limit {
		smallHeap = smallHeap[:limit]
	}
	return
}
