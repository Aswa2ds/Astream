package astream

import (
	"sort"
	"testing"
)

func BenchmarkStream(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var list []interface{}
		for i := 0; i < 10000; i++ {
			list = append(list, 10000-i)
		}
		Stream(list).Sort(func(a, b interface{}) bool {
			return a.(int) < b.(int)
		}).Collect()
	}
}

type itfs []interface{}

func (itfs *itfs) Len() int {
	return len(*itfs)
}

func (itfs *itfs) Less(i, j int) bool {
	return (*itfs)[i].(int) > (*itfs)[j].(int)
}

func (itfs *itfs) Swap(i, j int) {
	(*itfs)[i], (*itfs)[j] = (*itfs)[i], (*itfs)[i]
}

func BenchmarkNormal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var list itfs
		for i := 0; i < 10000; i++ {
			list = append(list, 10000-i)
		}
		//for j := range list {
		//	list[j] = list[j].(int) + 10
		//}
		sort.Sort(&list)
		//for j := 0; j < len(list); j++ {
		//	if list[j].(int) <= 1000 || list[j].(int) >= 5000 {
		//		list = append(list[:j], list[j+1:]...)
		//	}
		//}
		//sort.Sort(&list)
		//for j := range list {
		//	_ = list[j] + 10
		//}

	}
}

//func BenchmarkInterfaceCast(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		var list []int
//		for j := 0; j < 100000; j++ {
//			list = append(list, i)
//		}
//		var a interface{}
//		for _, j := range list {
//			a = j
//		}
//		_ = a
//	}
//}
//
//func BenchmarkDoNotCast(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		var list []int
//		for j := 0; j < 100000; j++ {
//			list = append(list, i)
//		}
//		var a int
//		for _, j := range list {
//			a = j
//		}
//		_ = a
//	}
//}
