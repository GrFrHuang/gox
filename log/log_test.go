package log

import (
	"testing"
	"github.com/tealeg/xlsx"
	"fmt"
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

func TestExcel(t *testing.T) {
	file, err := xlsx.OpenFile("/home/huang/film/黄泽元-绩效考核表.xlsx")
	fmt.Println(file.Sheet["Sheet1"].Rows[4].Cells[2].Value, err)
}
