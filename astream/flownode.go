package astream

import "sync"

type flowNode struct {
	channel     chan interface{}
	chSize      int
	limitOrSkip int
	reduceBase  interface{}
	heap        asHeap
	operator    interface{}
	nodeType    NodeType
	next        *flowNode
}

type handleFunc func(node *flowNode, resultChan flowResultChan)

var handleStream handleFunc = func(node *flowNode, resultChan flowResultChan) {
	output := node.next.channel
	defer close(output)
	for _, num := range node.heap.data {
		output <- num
	}
}

var handleForEach handleFunc = func(node *flowNode, resultChan flowResultChan) {
	defer close(resultChan)
	forEachFunc := node.operator.(ForEachFunc)
	for value, ok := <-node.channel; ok; value, ok = <-node.channel {
		forEachFunc(value)
	}
}

var handleMap handleFunc = func(node *flowNode, resultChan flowResultChan) {
	output := node.next.channel
	defer close(output)
	mapFunc := node.operator.(MapFunc)
	for value, ok := <-node.channel; ok; value, ok = <-node.channel {
		output <- mapFunc(value)
	}
}

var handleFlatMap handleFunc = func(node *flowNode, resultChan flowResultChan) {
	output := node.next.channel
	defer close(output)
	flatMapFunc := node.operator.(FlatMapFunc)
	waitGroup := sync.WaitGroup{}
	for value, ok := <-node.channel; ok; value, ok = <-node.channel {
		subFlow := flatMapFunc(value)
		waitGroup.Add(1)
		go func() {
			for _, datum := range (*subFlow)[0].heap.data {
				output <- datum
			}
			waitGroup.Done()
		}()
	}
	waitGroup.Wait()
}

var handleReduce handleFunc = func(node *flowNode, resultChan flowResultChan) {
	defer close(resultChan)
	result := node.reduceBase
	reduceFunc := node.operator.(ReduceFunc)
	for value, ok := <-node.channel; ok; value, ok = <-node.channel {
		result = reduceFunc(result, value)
	}
	resultChan <- result
}

var handleFilter handleFunc = func(node *flowNode, resultChan flowResultChan) {
	output := node.next.channel
	defer close(output)
	filterFunc := node.operator.(FilterFunc)
	for value, ok := <-node.channel; ok; value, ok = <-node.channel {
		if filterFunc(value) {
			output <- value
		}
	}
}

var handleSort handleFunc = func(node *flowNode, resultChan flowResultChan) {
	output := node.next.channel
	defer close(output)
	for value, ok := <-node.channel; ok; value, ok = <-node.channel {
		node.heap.push(value)
	}
	for !node.heap.isEmpty() {
		output <- node.heap.pop()
	}
}

var handleDistinct handleFunc = func(node *flowNode, resultChan flowResultChan) {
	output := node.next.channel
	defer close(output)
	dict := make(map[interface{}]bool)
	for value, ok := <-node.channel; ok; value, ok = <-node.channel {
		if _, ok := dict[value]; !ok {
			dict[value] = true
			output <- value
		}
	}
}

var handleLimit handleFunc = func(node *flowNode, resultChan flowResultChan) {
	output := node.next.channel
	defer close(output)
	limit := node.limitOrSkip
	for value, ok := <-node.channel; ok; value, ok = <-node.channel {
		if limit <= 0 {
			break
		}
		output <- value
		limit--
	}
	for _, ok := <-node.channel; ok; _, ok = <-node.channel {
	}
}

var handleSkip handleFunc = func(node *flowNode, resultChan flowResultChan) {
	output := node.next.channel
	defer close(output)
	skip := node.limitOrSkip
	for _, ok := <-node.channel; ok; _, ok = <-node.channel {
		skip--
		if skip <= 0 {
			break
		}
	}
	for value, ok := <-node.channel; ok; value, ok = <-node.channel {
		output <- value
	}
}

var handleAllMatch handleFunc = func(node *flowNode, resultChan flowResultChan) {
	defer close(resultChan)
	matchFunc := node.operator.(MatchFunc)
	flag := true
	for value, ok := <-node.channel; ok; value, ok = <-node.channel {
		if !matchFunc(value) {
			flag = false
			resultChan <- flag
			break
		}
	}
	if flag {
		resultChan <- flag
	}
	for _, ok := <-node.channel; ok; _, ok = <-node.channel {
	}
}

var handleAnyMatch handleFunc = func(node *flowNode, resultChan flowResultChan) {
	defer close(resultChan)
	matchFunc := node.operator.(MatchFunc)
	flag := false
	for value, ok := <-node.channel; ok; value, ok = <-node.channel {
		if matchFunc(value) {
			flag = true
			resultChan <- flag
			break
		}
	}
	if !flag {
		resultChan <- flag
	}
	for _, ok := <-node.channel; ok; _, ok = <-node.channel {
	}
}

var handleNoneMatch handleFunc = func(node *flowNode, resultChan flowResultChan) {
	defer close(resultChan)
	matchFunc := node.operator.(MatchFunc)
	flag := true
	for value, ok := <-node.channel; ok; value, ok = <-node.channel {
		if matchFunc(value) {
			flag = false
			resultChan <- flag
			break
		}
	}
	if flag {
		resultChan <- flag
	}
	for _, ok := <-node.channel; ok; _, ok = <-node.channel {
	}
}

var handleEmpty handleFunc = func(node *flowNode, resultChan flowResultChan) {
	defer close(resultChan)
	if _, ok := <-node.channel; !ok {
		resultChan <- true
	} else {
		resultChan <- false
	}
	for _, ok := <-node.channel; ok; _, ok = <-node.channel {
	}
}

var handleCollect handleFunc = func(node *flowNode, resultChan flowResultChan) {
	defer close(resultChan)
	flowResult := make([]interface{}, 0)
	for value, ok := <-node.channel; ok; value, ok = <-node.channel {
		flowResult = append(flowResult, value)
	}
	resultChan <- flowResult
}

var handleCount handleFunc = func(node *flowNode, resultChan flowResultChan) {
	defer close(resultChan)
	var count int
	for _, ok := <-node.channel; ok; _, ok = <-node.channel {
		count++
	}
	resultChan <- count
}

var handleMax handleFunc = func(node *flowNode, resultChan flowResultChan) {
	defer close(resultChan)
	var max interface{}
	if value, ok := <-node.channel; ok {
		max = value
	} else {
		resultChan <- nil
	}
	maxFunc := node.operator.(MaxFunc)
	for value, ok := <-node.channel; ok; value, ok = <-node.channel {
		max = maxFunc(max, value)
	}
	resultChan <- max
}

var handleMin handleFunc = func(node *flowNode, resultChan flowResultChan) {
	defer close(resultChan)
	var max interface{}
	if value, ok := <-node.channel; ok {
		max = value
	} else {
		resultChan <- nil
	}
	minFunc := node.operator.(MinFunc)
	for value, ok := <-node.channel; ok; value, ok = <-node.channel {
		max = minFunc(max, value)
	}
	resultChan <- max
}
