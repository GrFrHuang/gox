package md5

import (
	"testing"
	"github.com/GrFrHuang/gox/log"
)

type Student struct {
	Name string
	Age  int
}

func TestOnceMD5(t *testing.T) {
	s := Student{
		Name: "GrFrHuang",
		Age:  24,
	}
	res, err := OnceMD5(s)
	log.Debug(res, err)
}
