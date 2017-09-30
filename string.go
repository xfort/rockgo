package rockgo

import (
	"crypto/sha1"
	"encoding/hex"
)

func Substr(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		return ""
	}

	if end < 0 || end > length {
		return ""
	}
	return string(rs[start:end])

}

func Md5Hash(str string) string {
	md5HashEr := sha1.New()
	md5HashEr.Write([]byte(str))
	return hex.EncodeToString(md5HashEr.Sum(nil))
}
