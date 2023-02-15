package encrypts

import (
	"crypto/md5"
	"fmt"
	"io"
)

func Md5(str string) string {
	hash := md5.New()
	_, _ = io.WriteString(hash, str)
	return fmt.Sprintf("%x", hash.Sum(nil))
}
