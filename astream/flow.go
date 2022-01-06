package astream

type flowInterface interface {
	ForEach(eachFunc ForEachFunc)
	Map(mapFunc MapFunc) *Flow
	Reduce(reduceBase interface{}, reduceFunc ReduceFunc) FlowResult
	Filter(filterFunc FilterFunc) *Flow
	Sort(sortFunc SortFunc) *Flow
	Distinct() *Flow
	Limit(n int) *Flow
	Skip(n int) *Flow
	AllMatch(matchFunc MatchFunc) bool
	AnyMatch(matchFunc MatchFunc) bool
	NoneMatch(matchFunc MatchFunc) bool
	Empty() bool
	Collect() FlowResult
	Count() int
}

type Flow []*flowNode

func (flow *Flow) ForEach(forEachFunc ForEachFunc) {
	node := &flowNode{
		operator: forEachFunc,
		nodeType: opForeach,
		channel:  flow.genChannel(),
	}
	flow.append(node)
	flow.run()
}

func (flow *Flow) Map(mapFunc MapFunc) *Flow {
	node := &flowNode{
		operator: mapFunc,
		nodeType: opMap,
		channel:  flow.genChannel(),
	}
	flow.append(node)
	return flow
}

func (flow *Flow) FlatMap(flatMapFunc FlatMapFunc) *Flow {
	node := &flowNode{
		operator: flatMapFunc,
		nodeType: opFlatMap,
		channel:  flow.genChannel(),
	}
	flow.append(node)
	return flow
}

func (flow *Flow) Reduce(reduceBase interface{}, reduceFunc ReduceFunc) FlowResult {
	node := &flowNode{
		operator:   reduceFunc,
		nodeType:   opReduce,
		reduceBase: reduceBase,
		channel:    flow.genChannel(),
	}
	flow.append(node)
	return flow.run()
}

func (flow *Flow) Filter(filterFunc FilterFunc) *Flow {
	node := &flowNode{
		operator: filterFunc,
		nodeType: opFilter,
		channel:  flow.genChannel(),
	}
	flow.append(node)
	return flow
}

func (flow *Flow) Sort(sortFunc SortFunc) *Flow {
	node := &flowNode{
		channel: flow.genChannel(),
		heap: asHeap{
			data:     make([]interface{}, 0),
			SortFunc: sortFunc,
		},
		operator: sortFunc,
		nodeType: "opSort",
	}
	flow.append(node)
	return flow
}

func (flow *Flow) Distinct() *Flow {
	node := &flowNode{
		channel:  flow.genChannel(),
		nodeType: opDistinct,
	}
	flow.append(node)
	return flow
}

func (flow *Flow) Limit(n int) *Flow {
	node := &flowNode{
		channel:     flow.genChannel(),
		nodeType:    opLimit,
		limitOrSkip: n,
	}
	flow.append(node)
	return flow
}

func (flow *Flow) Skip(n int) *Flow {
	node := &flowNode{
		channel:     flow.genChannel(),
		nodeType:    opSkip,
		limitOrSkip: n,
	}
	flow.append(node)
	return flow
}

func (flow *Flow) AllMatch(matchFunc MatchFunc) bool {
	node := &flowNode{
		channel:  flow.genChannel(),
		nodeType: opAllMatch,
		operator: matchFunc,
	}
	flow.append(node)
	return flow.run().(bool)
}

func (flow *Flow) AnyMatch(matchFunc MatchFunc) bool {
	node := &flowNode{
		channel:  flow.genChannel(),
		nodeType: opAnyMatch,
		operator: matchFunc,
	}
	flow.append(node)
	return flow.run().(bool)
}

func (flow *Flow) NoneMatch(matchFunc MatchFunc) bool {
	node := &flowNode{
		channel:  flow.genChannel(),
		nodeType: opNoneMatch,
		operator: matchFunc,
	}
	flow.append(node)
	return flow.run().(bool)
}

func (flow *Flow) Empty() bool {
	node := &flowNode{
		channel:  flow.genChannel(),
		nodeType: opEmpty,
	}
	flow.append(node)
	return flow.run().(bool)
}

func (flow *Flow) Collect() FlowResult {
	node := &flowNode{
		channel:  flow.genChannel(),
		nodeType: opCollect,
	}
	flow.append(node)
	return flow.run()
}

func (flow *Flow) Count() int {
	node := &flowNode{
		channel:  flow.genChannel(),
		nodeType: opCount,
	}
	flow.append(node)
	return flow.run().(int)
}

func (flow *Flow) Max(maxFunc MaxFunc) FlowResult {
	node := &flowNode{
		channel:  flow.genChannel(),
		nodeType: opMax,
		operator: maxFunc,
	}
	flow.append(node)
	return flow.run()
}

func (flow *Flow) Min(minFunc MinFunc) FlowResult {
	node := &flowNode{
		channel:  flow.genChannel(),
		nodeType: opMin,
		operator: minFunc,
	}
	flow.append(node)
	return flow.run()
}

func (flow *Flow) append(node *flowNode) {
	l := len(*flow)
	(*flow)[l-1].next = node
	*flow = append(*flow, node)
}

func (flow *Flow) run() FlowResult {
	resultChan := make(FlowResultChan)

	for _, flowNode := range *flow {
		handler := handleFuncMap[flowNode.nodeType]
		go handler(flowNode, resultChan)
	}
	flowResult := <-resultChan
	return flowResult
}

func (flow *Flow) genChannel() chan interface{} {
	return make(chan interface{}, (*flow)[0].chSize)
}
