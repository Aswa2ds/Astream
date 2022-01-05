package astream

type flowNode struct {
	ch          chan interface{}
	chSize      int
	limitOrSkip int
	heap        asHeap
	operator    interface{}
	nodeType    NodeType
	next        *flowNode
}

func handleStream(node *flowNode, resultChan FlowResultChan) {
	go func() {
		output := node.next.ch
		defer close(output)
		for _, num := range node.heap.data {
			output <- num
		}
	}()
}

func handleForEach(node *flowNode, resultChan FlowResultChan) {
	go func() {
		output := node.next.ch
		defer close(output)
		forEachFunc := node.operator.(ForEachFunc)
		for value, ok := <-node.ch; ok; value, ok = <-node.ch {
			forEachFunc(value)
		}
	}()
}

func handleMap(node *flowNode, resultChan FlowResultChan) {
	go func() {
		output := node.next.ch
		defer close(output)
		mapFunc := node.operator.(MapFunc)
		for value, ok := <-node.ch; ok; value, ok = <-node.ch {
			output <- mapFunc(value)
		}
	}()
}

func handleFilter(node *flowNode, resultChan FlowResultChan) {
	go func() {
		output := node.next.ch
		defer close(output)
		filterFunc := node.operator.(FilterFunc)
		for value, ok := <-node.ch; ok; value, ok = <-node.ch {
			if filterFunc(value) {
				output <- value
			}
		}
	}()
}

func handleSort(node *flowNode, resultChan FlowResultChan) {
	go func() {
		output := node.next.ch
		defer close(output)
		for value, ok := <-node.ch; ok; value, ok = <-node.ch {
			node.heap.push(value)
		}
		for !node.heap.isEmpty() {
			output <- node.heap.pop()
		}
	}()
}

func handleDistinct(node *flowNode, resultChan FlowResultChan) {
	go func() {
		output := node.next.ch
		defer close(output)
		dict := make(map[interface{}]bool)
		for value, ok := <-node.ch; ok; value, ok = <-node.ch {
			if _, ok := dict[value]; !ok {
				dict[value] = true
				output <- value
			}
		}
	}()
}

func handleLimit(node *flowNode, resultChan FlowResultChan) {
	go func() {
		output := node.next.ch
		defer close(output)
		limit := node.limitOrSkip
		for value, ok := <-node.ch; ok; value, ok = <-node.ch {
			if limit <= 0 {
				break
			}
			output <- value
			limit--
		}
		for _, ok := <-node.ch; ok; _, ok = <-node.ch {
		}
	}()
}

func handleSkip(node *flowNode, resultChan FlowResultChan) {
	go func() {
		output := node.next.ch
		defer close(output)
		skip := node.limitOrSkip
		for _, ok := <-node.ch; ok; _, ok = <-node.ch {
			skip--
			if skip <= 0 {
				break
			}
		}
		for value, ok := <-node.ch; ok; value, ok = <-node.ch {
			output <- value
		}
	}()
}

func handleAllMatch(node *flowNode, resultChan FlowResultChan) {
	go func() {
		output := node.next.ch
		defer close(output)
		matchFunc := node.operator.(MatchFunc)
		flag := true
		for value, ok := <-node.ch; ok; value, ok = <-node.ch {
			if !matchFunc(value) {
				flag = false
				output <- flag
				break
			}
		}
		if flag {
			output <- flag
		}
		for _, ok := <-node.ch; ok; _, ok = <-node.ch {
		}
	}()
}

func handleAnyMatch(node *flowNode, resultChan FlowResultChan) {
	go func() {
		output := node.next.ch
		defer close(output)
		matchFunc := node.operator.(MatchFunc)
		flag := false
		for value, ok := <-node.ch; ok; value, ok = <-node.ch {
			if matchFunc(value) {
				flag = true
				output <- flag
				break
			}
		}
		if !flag {
			output <- flag
		}
		for _, ok := <-node.ch; ok; _, ok = <-node.ch {
		}
	}()
}

func handleNoneMatch(node *flowNode, resultChan FlowResultChan) {
	go func() {
		output := node.next.ch
		defer close(output)
		matchFunc := node.operator.(MatchFunc)
		flag := true
		for value, ok := <-node.ch; ok; value, ok = <-node.ch {
			if matchFunc(value) {
				flag = false
				output <- flag
				break
			}
		}
		if flag {
			output <- flag
		}
		for _, ok := <-node.ch; ok; _, ok = <-node.ch {
		}
	}()
}

func handleCollect(node *flowNode, resultChan FlowResultChan) {
	go func() {
		flowResult := make([]interface{}, 0)
		for value, ok := <-node.ch; ok; value, ok = <-node.ch {
			flowResult = append(flowResult, value)
		}
		resultChan <- flowResult
	}()
}
