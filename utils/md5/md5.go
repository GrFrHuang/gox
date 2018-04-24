package md5

import (
	"github.com/GrFrHuang/gox/log"
	"crypto/md5"
	"fmt"
)

// Once md5 sum digest hash value.
// Use gob package.
func OnceMD5(content interface{}) (result string, err error) {
	hash := md5.New()
	value := fmt.Sprintf("%s", content)
	_, err = hash.Write([]byte(value))
	if err != nil {
		log.Error(err)
		return
	}
	// Transform hash byte value to be hex value.
	result = fmt.Sprintf("%x", hash.Sum(nil))
	return
}
