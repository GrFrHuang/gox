package main

import (
	"net/http"
	"fmt"
	"github.com/GrFrHuang/gox/tools/thirdpayment/models"
)

func main() {
	// 设置路由，如果访问/，则调用index方法
	http.HandleFunc("/v1/aliPayNotify", models.Index)
	var a = 0
	fmt.Println(a)
	a = 2
	// 启动web服务，监听9090端口
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}
