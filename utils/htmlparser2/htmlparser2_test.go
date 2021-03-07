package htmlparser2

import (
	"fmt"
	"testing"
)

// TestParseURL 测试生成URL预览
func TestParseURL(t *testing.T) {
	title, desc, img, err := ParseURL("https://detail.1688.com/offer/587028735680.html?spm=a2609.11209760.j661dm7m.6.44292de1umP81Z&tracelog=p4p&clickid=bbc6ba83b35a4657885c1d9f2766847c&sessionid=17cc5c787f511cef575df4749836baf1")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(title)
	fmt.Println(desc)
	fmt.Println(img)
}
