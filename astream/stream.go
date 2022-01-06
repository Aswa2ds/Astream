package astream

import (
	"math"
	"reflect"
)

func Stream(list interface{}) *Flow {
	data := toInterfaceSlice(list)
	flow := make(Flow, 0)
	node := &flowNode{
		heap: asHeap{
			data:     data,
			SortFunc: nil,
		},
		//chSize: len(list)/10,
		chSize:   int(math.Sqrt(float64(len(data)))),
		nodeType: opStream,
	}
	flow = append(flow, node)
	return &flow
}

func toInterfaceSlice(list interface{}) []interface{} {
	interfaceSlice := make([]interface{}, 0)
	if reflect.TypeOf(list).Kind() == reflect.Slice {
		s := reflect.ValueOf(list)
		for i := 0; i < s.Len(); i++ {
			interfaceSlice = append(interfaceSlice, s.Index(i).Interface())
		}
	} else {
		panic("Stream should be used on slice")
	}
	return interfaceSlice
}
