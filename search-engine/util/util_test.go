package util

import (
	"fmt"
	"strings"
	"testing"
)

func TestRemovePunctuation(t *testing.T) {
	origin_str := "th is is !!,EGDsg .dfSFs gdfnij.///??.,[];'"
	str := RemovePunctuation(origin_str)
	str = strings.ToLower(str)
	str = RemoveSpace(str)
	fmt.Println(str)
}
