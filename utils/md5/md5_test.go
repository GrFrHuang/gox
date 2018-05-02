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
	res, err := OnceMD5(" game_key=aobohudong&secret_key=aobohudong&amount=0.01&order_no=1524652698666")
	log.Debug(res, err)
}
