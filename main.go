package main

import (
	"aswa2ds.cn/Astream/astream"
	"fmt"
)

func main() {
	var list []int
	for i := 0; i < 10000; i++ {
		list = append(list, i)
	}
	astream.Stream(astream.IntToInterface(list)).Map(func(a interface{}) interface{} {
		return a.(int) + 10
	}).Filter(func(a interface{}) bool {
		return a.(int) > 10
	}).Sort(func(a, b interface{}) bool {
		return a.(int) < b.(int)
	}).ForEach(func(a interface{}) {
		fmt.Println(a)
	}).Run()
}
