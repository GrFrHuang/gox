package pressure_valve

import (
	"net/http"
	"fmt"
	"github.com/GrFrHuang/gox/log"
	"testing"
)

var pv = NewPressureValve(10, 2000, 5, false)

func Handler(w http.ResponseWriter, r *http.Request) {
	err := pv.FlowFilter()
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Fprintln(w, "hello world")
}

func CreateErrHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintln(w, "hello world")
	pv.HaveError(fmt.Errorf("hava error"))
	fmt.Fprintln(w, "create error")
}

func TestPressureValve(t *testing.T) {
	http.HandleFunc("/index", Handler)
	http.HandleFunc("/err", CreateErrHandler)
	http.ListenAndServe("127.0.0.1:18001", nil)
	c := make(chan struct{})
	<-c
}
