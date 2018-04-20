package models

import (
	"net/http"
	"fmt"
)

func Index(w http.ResponseWriter, r *http.Request) {
	// 往w里写入内容，就会在浏览器里输出
	var b = 3
	var c = 4
	fmt.Fprintf(w, "ok", b, c)
}
