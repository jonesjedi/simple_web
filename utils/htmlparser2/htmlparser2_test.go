package htmlparser2

import (
	"fmt"
	"testing"
)

// TestParseURL 测试生成URL预览
func TestParseURL(t *testing.T) {
	title, desc, img, err := ParseURL("https://www.163.com")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(title)
	fmt.Println(desc)
	fmt.Println(img)
}
