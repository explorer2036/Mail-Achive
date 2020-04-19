package utils

import (
	"crypto/md5"
	"fmt"
)

// MD5Str - encode the string
func MD5Str(s string) string {
	d := []byte(s)
	b := md5.Sum(d)
	return fmt.Sprintf("%x", b)
}