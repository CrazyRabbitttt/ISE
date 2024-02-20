package util

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"regexp"
)

const (
	c1 = 0xcc9e2d51
	c2 = 0x1b873593
	c3 = 0x85ebca6b
	c4 = 0xc2b2ae35
	r1 = 15
	r2 = 13
	m  = 5
	n  = 0xe6546b64
)

var (
	Seed = uint32(1)
)

func Murmur3(key []byte) (hash uint32) {
	hash = Seed
	iByte := 0
	for ; iByte+4 <= len(key); iByte += 4 {
		k := uint32(key[iByte]) | uint32(key[iByte+1])<<8 | uint32(key[iByte+2])<<16 | uint32(key[iByte+3])<<24
		k *= c1
		k = (k << r1) | (k >> (32 - r1))
		k *= c2
		hash ^= k
		hash = (hash << r2) | (hash >> (32 - r2))
		hash = hash*m + n
	}

	var remainingBytes uint32
	switch len(key) - iByte {
	case 3:
		remainingBytes += uint32(key[iByte+2]) << 16
		fallthrough
	case 2:
		remainingBytes += uint32(key[iByte+1]) << 8
		fallthrough
	case 1:
		remainingBytes += uint32(key[iByte])
		remainingBytes *= c1
		remainingBytes = (remainingBytes << r1) | (remainingBytes >> (32 - r1))
		remainingBytes = remainingBytes * c2
		hash ^= remainingBytes
	}

	hash ^= uint32(len(key))
	hash ^= hash >> 16
	hash *= c3
	hash ^= hash >> 13
	hash *= c4
	hash ^= hash >> 16

	// 出发吧，狗嬷嬷！
	return
}

func RemovePunctuation(word string) string {
	reg := regexp.MustCompile(`[[:punct:]]`) // 定义要删除的标点符号范围为所有标点符号
	return reg.ReplaceAllString(word, "")    // 将匹配到的标点符号全部替换成空字符串
}

func RemoveSpace(word string) string {
	reg := regexp.MustCompile(`\s+`)
	return reg.ReplaceAllString(word, "")
}

func String2Int(str string) uint32 {
	// 获得一个 Hash 值
	return Murmur3([]byte(str))
}

func Uint32ToBytes(num uint32) []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, num)
	return buf
}

func Encoder(data interface{}) []byte {
	if data == nil {
		return nil
	}
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(data)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

// 从 Leveldb 中将读取到的value解析出来
func Decoder(data []byte, v interface{}) {
	if data == nil {
		return
	}
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(v)
	if err != nil {
		panic(err)
	}
}

func RemoveUint32ValueInArray(array []uint32, id uint32) []uint32 {
	var removeIndex int
	for index, value := range array {
		if value == id {
			removeIndex = index
			break
		}
	}
	return append(array[:removeIndex], array[removeIndex+1:]...)
}

func ExistInArrayUint32(array []uint32, id uint32) (int, bool) {
	for index, value := range array {
		if value == id {
			return index, true
		}
	}
	return -1, false
}

func ExistInArrayString(array []string, word string) (int, bool) {
	for index, value := range array {
		if word == value {
			return index, true
		}
	}
	return -1, false
}

// 计算最长公共子序列的长度
func CalculateLCS(text1 string, text2 string) int {
	m, n := len(text1), len(text2)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if text1[i-1] == text2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i][j-1], dp[i-1][j])
			}
		}
	}
	return dp[m][n]
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
