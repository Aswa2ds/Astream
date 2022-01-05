package astream

type flowNode struct {
	ch       chan interface{}
	chSize   int
	heap     asHeap
	operator interface{}
	nodeType NodeType
	next     *flowNode
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

func handleCollect(node *flowNode, resultChan FlowResultChan) {
	go func() {
		flowResult := make([]interface{}, 0)
		for value, ok := <-node.ch; ok; value, ok = <-node.ch {
			flowResult = append(flowResult, value)
		}
		resultChan <- flowResult
	}()
}
