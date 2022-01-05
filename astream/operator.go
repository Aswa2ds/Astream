package astream

const (
	STREAM  NodeType = "Stream"
	FOREACH NodeType = "FOREACH"
	MAP     NodeType = "Map"
	FILTER  NodeType = "FILTER"
	SORT    NodeType = "SORT"
	COLLECT NodeType = "COLLECT"
)

var handleFuncMap = map[NodeType]func(node *flowNode, resultChan FlowResultChan){
	STREAM:  handleStream,
	FOREACH: handleForEach,
	MAP:     handleMap,
	FILTER:  handleFilter,
	SORT:    handleSort,
	COLLECT: handleCollect,
}

type ForEachFunc func(a interface{})
type MapFunc func(a interface{}) interface{}
type FilterFunc func(a interface{}) bool
type SortFunc func(a, b interface{}) bool
