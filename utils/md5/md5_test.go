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
	//s := Student{
	//	Name: "GrFrHuang",
	//	Age:  24,
	//}
	//res, err := OnceMD5(s)
	res, err := OnceMD5("game_key=123123&secret_key=1231231&amount=0.01")
	log.Debug(res, err)
}
