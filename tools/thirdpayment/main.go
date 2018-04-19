package main

import (
	"net/http"
	"fmt"
)

func index(w http.ResponseWriter, r *http.Request) {
	// 往w里写入内容，就会在浏览器里输出
	fmt.Fprintf(w, "ok")
}

func main() {
	// 设置路由，如果访问/，则调用index方法
	http.HandleFunc("/v1/aliPayNotify", index)

	// 启动web服务，监听9090端口
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}

