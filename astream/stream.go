package astream

import "math"

type NodeType string

type Flow []*flowNode

type FlowResult interface{}
type FlowResultChan chan FlowResult

func StreamOf(list Interface) *Flow {
	flow := make(Flow, 0)
	node := &flowNode{
		heap: asHeap{
			data:     list.ToInterface(),
			SortFunc: nil,
		},
		nodeType: STREAM,
	}
	flow = append(flow, node)
	return &flow
}

func Stream(list []interface{}) *Flow {
	flow := make(Flow, 0)
	node := &flowNode{
		heap: asHeap{
			data:     list,
			SortFunc: nil,
		},
		//chSize: len(list)/10,
		chSize:   int(math.Sqrt(float64(len(list)))),
		nodeType: STREAM,
	}
	flow = append(flow, node)
	return &flow
}

func (flow *Flow) ForEach(forEachFunc ForEachFunc) *Flow {
	node := &flowNode{
		operator: forEachFunc,
		nodeType: FOREACH,
		ch:       make(chan interface{}, (*flow)[0].chSize),
	}
	flow.append(node)
	return flow
}

func (flow *Flow) Map(mapFunc MapFunc) *Flow {
	node := &flowNode{
		operator: mapFunc,
		nodeType: MAP,
		ch:       make(chan interface{}, (*flow)[0].chSize),
	}
	flow.append(node)
	return flow
}

func (flow *Flow) Filter(filterFunc FilterFunc) *Flow {
	node := &flowNode{
		operator: filterFunc,
		nodeType: FILTER,
		ch:       make(chan interface{}, (*flow)[0].chSize),
	}
	flow.append(node)
	return flow
}

func (flow *Flow) Sort(sortFunc SortFunc) *Flow {
	node := &flowNode{
		ch: make(chan interface{}, (*flow)[0].chSize),
		heap: asHeap{
			data:     make([]interface{}, 0),
			SortFunc: sortFunc,
		},
		operator: sortFunc,
		nodeType: "SORT",
	}
	flow.append(node)
	return flow
}

func (flow *Flow) append(node *flowNode) {
	l := len(*flow)
	(*flow)[l-1].next = node
	*flow = append(*flow, node)
}

func (flow *Flow) appendCollectNode() {
	node := &flowNode{
		ch:       make(chan interface{}),
		nodeType: COLLECT,
		next:     nil,
	}
	flow.append(node)
}

func (flow *Flow) Run() FlowResult {
	resultChan := make(FlowResultChan)

	flow.appendCollectNode()
	for _, flowNode := range *flow {
		handleFunc := handleFuncMap[flowNode.nodeType]
		handleFunc(flowNode, resultChan)
	}
	flowResult := <-resultChan
	return flowResult
}
