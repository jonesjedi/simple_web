package utils

import (
	"crypto/md5"
	"fmt"
)

func Md5Sum(src string) (dst string, err error) {

	data := []byte(src)
	has := md5.Sum(data)
	dst = fmt.Sprintf("%x", has) //将[]byte转成16进制

	return
}
