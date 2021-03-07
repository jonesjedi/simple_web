package htmlparser2

import (
	"fmt"
	"testing"
)

// TestParseURL 测试生成URL预览
func TestParseURL(t *testing.T) {
	title, desc, img, err := ParseURL("https://www.amazon.com/AmazonBasics-Accessory-Tray-6-Compartments/dp/B07FFTG66J/ref=sr_1_1_sspa?dchild=1&keywords=amazonbasics&pd_rd_r=6ece80d8-5878-4b90-bc2a-bad11bd770b2&pd_rd_w=xGfZN&pd_rd_wg=a1hh5&pf_rd_p=9349ffb9-3aaa-476f-8532-6a4a5c3da3e7&pf_rd_r=EPZ8DFEHNF3T7RVE4T78&qid=1614525108&sr=8-1-spons&psc=1&spLa=ZW5jcnlwdGVkUXVhbGlmaWVyPUE3M05NU1BNS1c0M0wmZW5jcnlwdGVkSWQ9QTA4ODEyMzMyNkJERVZTS1VZSVM1JmVuY3J5cHRlZEFkSWQ9QTAwOTY4MjJMN0hJUTEwOEFBSUkmd2lkZ2V0TmFtZT1zcF9hdGYmYWN0aW9uPWNsaWNrUmVkaXJlY3QmZG9Ob3RMb2dDbGljaz10cnVl")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(title)
	fmt.Println(desc)
	fmt.Println(img)
}
