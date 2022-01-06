package main

import (
	"aswa2ds.cn/Astream/astream"
	"fmt"
)

type class struct {
	name  string
	count int
}

type classes []class

func main() {
	//list := make([]int, 0)
	//rand.Seed(time.Now().UnixNano())
	//for i := 0; i < 10000; i++ {
	//	//list = append(list, rand.Intn(12423))
	//}
	//list := []string{"hello", "world", "cat", "dog", "banana"}
	list := classes{{
		name:  "1",
		count: 40,
	}, {
		name:  "1",
		count: 40,
	}}

	result := astream.Stream(list).Distinct().Collect()
	fmt.Println(result)
}
