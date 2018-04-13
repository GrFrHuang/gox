package md5

import (
	"crypto/md5"
	"io"
	"fmt"
	"github.com/GrFrHuang/gox/log"
	"bytes"
	"encoding/gob"
)

// Once md5 sum digest hash value.
// Use gob package.
func OnceMD5(content interface{}) (result string, err error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err = enc.Encode(content)
	if err != nil {
		log.Error(err)
		return
	}
	src := bytes.NewReader(buf.Bytes())
	hash := md5.New()
	// Type Hash extends from io.Writer.
	// Parent interface can be convert type to child interface.
	_, err = io.Copy(hash, src)
	if err != nil {
		log.Error(err)
		return
	}
	// Transform hash byte value to be hex value.
	result = fmt.Sprintf("%x", hash.Sum(nil))
	return
}
