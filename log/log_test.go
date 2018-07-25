package log

import (
	"testing"
)

type Student struct {
	Name string
	Age  int
}

func TestLog(t *testing.T) {
	stu := Student{
		Name: "GrFrHuang",
		Age:  24,
	}
	//logger := NewLogger()
	Warn("hello ", stu)
}
